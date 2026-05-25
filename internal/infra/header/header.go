package header

import (
	"net/http"

	"market-api-gateway/internal/core/dto"
)

// Injector implements port.HeaderInjector.
type Injector struct{}

func New() *Injector { return &Injector{} }

// InjectHeaders adds X-User-ID and X-User-Role to the upstream request.
func (h *Injector) InjectHeaders(req *http.Request, claims *dto.JWTClaims) {
	if claims == nil {
		return
	}
	req.Header.Set("X-User-ID", claims.Sub)
	req.Header.Set("X-User-Role", claims.Role)
}

// StripIncoming removes X-User-* headers from the incoming client request.
func (h *Injector) StripIncoming(req *http.Request) {
	req.Header.Del("X-User-ID")
	req.Header.Del("X-User-Role")
}
