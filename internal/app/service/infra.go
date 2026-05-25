package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"market-api-gateway/internal"
	"market-api-gateway/internal/core/domain"
	"market-api-gateway/internal/infra/auth"
	infraCORS "market-api-gateway/internal/infra/cors"
	"market-api-gateway/internal/infra/header"
	"market-api-gateway/internal/infra/outbound/httpclient"
	"market-api-gateway/internal/infra/proxy"
	"market-api-gateway/internal/infra/rate"
	"market-api-gateway/internal/platform/metrics"
	"market-api-gateway/internal/platform/tracing"
)

// Infra groups all shared infrastructure dependencies.
type Infra struct {
	Cfg             *internal.Config
	Log             *slog.Logger
	Metrics         *metrics.Metrics
	Auth            *auth.Authenticator
	Router          *proxy.Router
	Proxy           *proxy.Proxy
	RateLimiter     *rate.RateLimiter
	HeaderInjector  *header.Injector
	CORSMiddleware  func(http.Handler) http.Handler
	HTTPClients     map[domain.ServiceName]*httpclient.Client
	shutdownTracing tracing.ShutdownFunc
}

func initInfra(ctx context.Context, cfg *internal.Config, log *slog.Logger) (*Infra, error) {
	initCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Tracing
	shutdownTracing, err := tracing.Init(initCtx, tracing.Config{
		Enabled:     cfg.OTel.Enabled,
		Exporter:    cfg.OTel.Exporter,
		Endpoint:    cfg.OTel.Endpoint,
		ServiceName: cfg.OTel.ServiceName,
	})
	if err != nil {
		return nil, fmt.Errorf("tracing: %w", err)
	}

	// Metrics
	m := metrics.New()

	// HTTP clients (one per upstream)
	clients := map[domain.ServiceName]*httpclient.Client{
		domain.ServiceAuth: httpclient.New(httpclient.Config{
			Target:     "auth",
			BaseURL:    cfg.ServiceURLs.AuthURL,
			MaxRetries: 3,
			Timeout:    cfg.UpstreamTimeout,
		}, m, log),
		domain.ServiceCatalog: httpclient.New(httpclient.Config{
			Target:     "catalog",
			BaseURL:    cfg.ServiceURLs.CatalogURL,
			MaxRetries: 3,
			Timeout:    cfg.UpstreamTimeout,
		}, m, log),
		domain.ServiceCart: httpclient.New(httpclient.Config{
			Target:     "cart",
			BaseURL:    cfg.ServiceURLs.CartURL,
			MaxRetries: 3,
			Timeout:    cfg.UpstreamTimeout,
		}, m, log),
		domain.ServiceOrder: httpclient.New(httpclient.Config{
			Target:     "order",
			BaseURL:    cfg.ServiceURLs.OrderURL,
			MaxRetries: 3,
			Timeout:    cfg.UpstreamTimeout,
		}, m, log),
	}

	// Proxy
	p, err := proxy.NewProxy(proxy.ServiceClients{
		AuthURL:    cfg.ServiceURLs.AuthURL,
		CatalogURL: cfg.ServiceURLs.CatalogURL,
		CartURL:    cfg.ServiceURLs.CartURL,
		OrderURL:   cfg.ServiceURLs.OrderURL,
		Clients:    clients,
	}, m, log)
	if err != nil {
		_ = shutdownTracing(ctx)
		return nil, fmt.Errorf("proxy: %w", err)
	}

	// Auth (JWKS)
	authenticator, err := auth.New(ctx, cfg.ServiceURLs.AuthURL, cfg.JWKS.RefreshInterval, log)
	if err != nil {
		_ = shutdownTracing(ctx)
		return nil, fmt.Errorf("auth: %w", err)
	}

	return &Infra{
		Cfg:            cfg,
		Log:            log,
		Metrics:        m,
		Auth:           authenticator,
		Router:         proxy.NewRouter(),
		Proxy:          p,
		RateLimiter:    rate.New(cfg.RateLimit.RPS),
		HeaderInjector: header.New(),
		CORSMiddleware: infraCORS.New(cfg.CORS.AllowedOrigins),
		HTTPClients:    clients,
		shutdownTracing: shutdownTracing,
	}, nil
}

func (i *Infra) Shutdown(ctx context.Context) {
	_ = i.shutdownTracing(ctx)
}

// upstreamTimeout returns a reasonable timeout for the graceful shutdown HTTP server.
func (i *Infra) upstreamTimeout() time.Duration {
	return i.Cfg.UpstreamTimeout
}
