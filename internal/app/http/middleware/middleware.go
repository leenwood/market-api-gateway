package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"

	"market-api-gateway/internal/core/dto"
	"market-api-gateway/internal/core/port"
	"market-api-gateway/internal/core/service"
	"market-api-gateway/internal/platform/logger"
	"market-api-gateway/internal/platform/metrics"
)

type ctxKey int

const claimsKey ctxKey = iota

// Chain composes middleware handlers (outermost first).
func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// ClaimsFromContext retrieves JWT claims stored in the request context.
func ClaimsFromContext(ctx context.Context) *dto.JWTClaims {
	v, _ := ctx.Value(claimsKey).(*dto.JWTClaims)
	return v
}

// Recover catches panics and returns 500.
func Recover(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.ErrorContext(r.Context(), "panic recovered", "panic", rec, "stack", string(debug.Stack()))
					writeError(w, http.StatusInternalServerError, "internal server error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// Logger logs each request with method, path, status, and latency; records Prometheus RED metrics.
func Logger(log *slog.Logger, m *metrics.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Extract trace_id from the active OTel span.
			spanCtx := trace.SpanFromContext(r.Context()).SpanContext()
			ctx := r.Context()
			if spanCtx.IsValid() {
				ctx = logger.WithTraceID(ctx, spanCtx.TraceID().String())
			}
			r = r.WithContext(ctx)

			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rw, r)

			dur := time.Since(start)
			l := logger.FromContext(r.Context(), log)
			l.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.status,
				"duration_ms", dur.Milliseconds(),
			)

			status := strconv.Itoa(rw.status)
			m.HTTPRequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
			m.HTTPRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(dur.Seconds())
		})
	}
}

// RequestID reads X-Request-ID or generates a UUID; writes it back to the response.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = uuid.NewString()
		}
		ctx := logger.WithRequestID(r.Context(), id)
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// MaxBodySize limits the request body to maxBytes.
func MaxBodySize(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}

// HeaderStrip removes X-User-* headers from incoming client requests.
func HeaderStrip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Del("X-User-ID")
		r.Header.Del("X-User-Role")
		next.ServeHTTP(w, r)
	})
}

// JWTAuth extracts and validates the Bearer token; attaches claims to context.
// Returns 401 on missing/invalid token, wraps service errors into HTTP codes.
func JWTAuth(auth port.Authenticator, log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractBearer(r)
			if token == "" {
				writeError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			claims, err := auth.ValidateToken(r.Context(), token)
			if err != nil {
				logger.FromContext(r.Context(), log).Warn("jwt validation failed", "err", err)
				writeError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RateLimit enforces per-IP rate limits; returns 429 on denial.
func RateLimit(rl port.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIP(r)
			if !rl.Allow(ip) {
				w.Header().Set("Retry-After", "1")
				writeError(w, http.StatusTooManyRequests, "rate limit exceeded")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// CartOwnership ensures the userID in the cart path matches the authenticated user.
func CartOwnership(router port.Router, log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			route, params, ok := router.GetRoute(r.URL.Path, r.Method)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			claims := ClaimsFromContext(r.Context())
			if err := service.CheckAccess(route, claims, params); err != nil {
				switch {
				case fmt.Sprintf("%v", err) == service.ErrForbidden.Error():
					writeError(w, http.StatusForbidden, "access denied")
				default:
					writeError(w, http.StatusUnauthorized, "unauthorized")
				}
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// --- helpers ---

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func extractBearer(r *http.Request) string {
	v := r.Header.Get("Authorization")
	if after, ok := strings.CutPrefix(v, "Bearer "); ok {
		return after
	}
	return ""
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.SplitN(xff, ",", 2)[0]
	}
	return r.RemoteAddr
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
