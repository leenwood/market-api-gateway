package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"market-api-gateway/internal"
	"market-api-gateway/internal/app/adapter"
	"market-api-gateway/internal/app/http/handler"
	"market-api-gateway/internal/app/http/middleware"
	"market-api-gateway/internal/core/domain"
	"market-api-gateway/internal/core/port"
	coreservice "market-api-gateway/internal/core/service"
	"market-api-gateway/internal/platform/logger"
)

const (
	shutdownTimeout = 30 * time.Second
	maxBodySize     = 1 << 20 // 1 MiB
)

func RunServer(ctx context.Context) error {
	cfg, err := internal.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	log := logger.New(cfg.Log.Level, cfg.Log.Format)
	log.Info("starting api-gateway", "addr", cfg.HTTP.Addr)

	infra, err := initInfra(ctx, cfg, log)
	if err != nil {
		return fmt.Errorf("infra init: %w", err)
	}

	srv := buildServer(cfg, infra, log)

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()
	log.Info("server listening", "addr", cfg.HTTP.Addr)

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		log.Info("shutting down server")
	}

	shutCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutCtx); err != nil {
		log.Error("server shutdown error", "err", err)
	}
	infra.Shutdown(shutCtx)
	log.Info("server stopped")
	return nil
}

func buildServer(cfg *internal.Config, infra *Infra, log *slog.Logger) *http.Server {
	mux := http.NewServeMux()

	// Static gateway endpoints — no auth, no rate limit.
	mux.HandleFunc("GET /health", handler.HealthHandler)
	mux.Handle("GET /metrics", infra.Metrics.Handler())
	mux.Handle("GET /ready", handler.ReadyHandler(buildPingers(cfg)))
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// All other requests — proxy handler.
	mux.Handle("/", buildProxyHandler(infra, log))

	// Middleware chain (outermost → innermost).
	otelFilter := otelhttp.WithFilter(func(r *http.Request) bool {
		return r.URL.Path != "/metrics" && r.URL.Path != "/health"
	})
	chain := middleware.Chain(
		mux,
		func(h http.Handler) http.Handler {
			return otelhttp.NewHandler(h, "api-gateway", otelFilter)
		},
		middleware.Recover(log),
		middleware.Logger(log, infra.Metrics),
		middleware.RequestID,
		middleware.MaxBodySize(maxBodySize),
	)

	return &http.Server{
		Addr:         cfg.HTTP.Addr,
		Handler:      chain,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}
}

func buildProxyHandler(infra *Infra, log *slog.Logger) http.Handler {
	core := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route, params, ok := infra.Router.GetRoute(r.URL.Path, r.Method)
		if !ok {
			http.NotFound(w, r)
			return
		}

		if route.Type == domain.Protected {
			// Validate JWT and store claims in context via JWTAuth middleware, then proceed.
			middleware.JWTAuth(infra.Auth, log)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				claims := middleware.ClaimsFromContext(r.Context())
				if err := coreservice.CheckAccess(route, claims, params); err != nil {
					switch {
					case errors.Is(err, coreservice.ErrForbidden):
						writeJSON(w, http.StatusForbidden, "access denied")
					default:
						writeJSON(w, http.StatusUnauthorized, "unauthorized")
					}
					return
				}
				infra.HeaderInjector.InjectHeaders(r, claims)
				infra.Proxy.ServeHTTP(w, r, route.Service)
			})).ServeHTTP(w, r)
			return
		}

		infra.Proxy.ServeHTTP(w, r, route.Service)
	})

	return infra.CORSMiddleware(
		middleware.HeaderStrip(
			middleware.RateLimit(infra.RateLimiter)(core),
		),
	)
}

func buildPingers(cfg *internal.Config) []port.Pinger {
	client := &http.Client{Timeout: 3 * time.Second}
	return []port.Pinger{
		adapter.NewHTTPPinger("auth", cfg.ServiceURLs.AuthURL, client),
		adapter.NewHTTPPinger("catalog", cfg.ServiceURLs.CatalogURL, client),
		adapter.NewHTTPPinger("cart", cfg.ServiceURLs.CartURL, client),
		adapter.NewHTTPPinger("order", cfg.ServiceURLs.OrderURL, client),
	}
}

func writeJSON(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = fmt.Fprintf(w, `{"error":%q}`, msg)
}
