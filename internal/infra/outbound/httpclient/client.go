package httpclient

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	"market-api-gateway/internal/platform/metrics"
)

const (
	cbMaxFailures  = 5
	cbOpenTimeout  = 30 * time.Second
	cbHalfProbes   = 2
	maxBackoffBase = 100 * time.Millisecond
	maxBackoffCap  = 30 * time.Second
)

type cbState int

const (
	cbClosed   cbState = iota
	cbOpen     cbState = iota
	cbHalfOpen cbState = iota
)

type Config struct {
	Target     string // short label for metrics/logs, e.g. "auth"
	BaseURL    string
	MaxRetries int
	Timeout    time.Duration
}

type Response struct {
	StatusCode int
	Body       []byte
	Header     http.Header
}

// Client is a resilient HTTP client with retry, circuit breaker, OTel tracing, and Prometheus metrics.
// It also implements http.RoundTripper for use as httputil.ReverseProxy.Transport.
type Client struct {
	cfg     Config
	inner   *http.Client
	metrics *metrics.Metrics
	log     *slog.Logger
	tracer  trace.Tracer

	cbState    cbState
	cbFailures int
	cbOpenedAt time.Time
	cbProbes   int
}

func New(cfg Config, m *metrics.Metrics, log *slog.Logger) *Client {
	if cfg.MaxRetries == 0 {
		cfg.MaxRetries = 3
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}
	return &Client{
		cfg:     cfg,
		inner:   &http.Client{Timeout: cfg.Timeout},
		metrics: m,
		log:     log,
		tracer:  otel.Tracer("httpclient/" + cfg.Target),
	}
}

// Do executes a request with retry and circuit breaker logic.
func (c *Client) Do(ctx context.Context, method, path string, body io.Reader) (*Response, error) {
	url := c.cfg.BaseURL + path

	ctx, span := c.tracer.Start(ctx, method+" "+path,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			semconv.HTTPRequestMethodKey.String(method),
			attribute.String("http.url", url),
		),
	)
	defer span.End()

	var lastErr error
	for attempt := 0; attempt <= c.cfg.MaxRetries; attempt++ {
		if attempt > 0 {
			delay := jitter(time.Duration(float64(maxBackoffBase) * float64(int(1)<<attempt)))
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		if !c.cbAllow() {
			return nil, errors.New("circuit breaker open")
		}

		req, err := http.NewRequestWithContext(ctx, method, url, body)
		if err != nil {
			return nil, fmt.Errorf("build request: %w", err)
		}
		// Inject W3C trace context.
		otel.GetTextMapPropagator().Inject(ctx, propagationCarrier(req.Header))

		resp, err := c.inner.Do(req)
		statusLabel := "error"
		if err == nil {
			statusLabel = strconv.Itoa(resp.StatusCode)
		}
		c.metrics.HTTPOutboundRequestTotal.WithLabelValues(c.cfg.Target, method, statusLabel).Inc()

		if err != nil {
			c.cbOnFailure()
			lastErr = err
			continue
		}

		if isRetryable(resp.StatusCode) {
			c.cbOnFailure()
			resp.Body.Close()
			lastErr = fmt.Errorf("upstream %d", resp.StatusCode)
			continue
		}

		c.cbOnSuccess()
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))
		if resp.StatusCode >= 400 {
			span.SetStatus(codes.Error, resp.Status)
		}

		return &Response{
			StatusCode: resp.StatusCode,
			Body:       respBody,
			Header:     resp.Header,
		}, nil
	}
	span.SetStatus(codes.Error, "all retries exhausted")
	return nil, fmt.Errorf("all retries exhausted: %w", lastErr)
}

// RoundTrip implements http.RoundTripper — used as ReverseProxy.Transport.
// It wraps the inner client transport with circuit breaker state and metrics,
// but does NOT retry (the proxy handles its own response).
func (c *Client) RoundTrip(req *http.Request) (*http.Response, error) {
	if !c.cbAllow() {
		return nil, errors.New("circuit breaker open")
	}

	otel.GetTextMapPropagator().Inject(req.Context(), propagationCarrier(req.Header))

	resp, err := c.inner.Do(req)
	statusLabel := "error"
	if err == nil {
		statusLabel = strconv.Itoa(resp.StatusCode)
	}
	c.metrics.HTTPOutboundRequestTotal.WithLabelValues(c.cfg.Target, req.Method, statusLabel).Inc()

	if err != nil {
		c.cbOnFailure()
		return nil, err
	}
	if resp.StatusCode >= 500 {
		c.cbOnFailure()
	} else {
		c.cbOnSuccess()
	}
	return resp, nil
}

// --- Circuit breaker ---

func (c *Client) cbAllow() bool {
	switch c.cbState {
	case cbOpen:
		if time.Since(c.cbOpenedAt) >= cbOpenTimeout {
			c.cbState = cbHalfOpen
			c.cbProbes = 0
			c.log.Warn("circuit breaker half-open", "target", c.cfg.Target)
			return true
		}
		return false
	default:
		return true
	}
}

func (c *Client) cbOnSuccess() {
	if c.cbState == cbHalfOpen {
		c.cbProbes++
		if c.cbProbes >= cbHalfProbes {
			c.cbState = cbClosed
			c.cbFailures = 0
			c.log.Warn("circuit breaker closed", "target", c.cfg.Target)
		}
		return
	}
	c.cbFailures = 0
}

func (c *Client) cbOnFailure() {
	c.cbFailures++
	if c.cbState == cbHalfOpen || c.cbFailures >= cbMaxFailures {
		c.cbState = cbOpen
		c.cbOpenedAt = time.Now()
		c.log.Warn("circuit breaker open", "target", c.cfg.Target, "failures", c.cbFailures)
	}
}

// --- Helpers ---

func isRetryable(code int) bool {
	switch code {
	case 429, 502, 503, 504:
		return true
	}
	return code >= 500
}

func jitter(base time.Duration) time.Duration {
	if base > maxBackoffCap {
		base = maxBackoffCap
	}
	// full jitter: rand(0, base)
	return time.Duration(rand.Int64N(int64(base) + 1))
}

// propagationCarrier adapts http.Header to otel TextMapCarrier.
type propagationCarrier http.Header

func (c propagationCarrier) Get(key string) string        { return http.Header(c).Get(key) }
func (c propagationCarrier) Set(key, val string)          { http.Header(c).Set(key, val) }
func (c propagationCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}
