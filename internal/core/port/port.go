package port

import (
	"context"
	"net/http"

	"market-api-gateway/internal/core/domain"
	"market-api-gateway/internal/core/dto"
)

// Router matches an incoming request path+method against the routing table.
type Router interface {
	GetRoute(path, method string) (*domain.Route, map[string]string, bool)
}

// Authenticator validates a Bearer token and returns the verified claims.
type Authenticator interface {
	ValidateToken(ctx context.Context, token string) (*dto.JWTClaims, error)
}

// RateLimiter enforces per-IP request rate limits.
type RateLimiter interface {
	Allow(ip string) bool
}

// HeaderInjector injects user context headers into an upstream request.
type HeaderInjector interface {
	InjectHeaders(req *http.Request, claims *dto.JWTClaims)
	StripIncoming(req *http.Request)
}

// Pinger checks connectivity to a named upstream service.
type Pinger interface {
	Ping(ctx context.Context) error
	Name() string
}
