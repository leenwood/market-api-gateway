package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"

	"market-api-gateway/internal/core/dto"
)

var (
	ErrTokenExpired = errors.New("token expired")
	ErrTokenInvalid = errors.New("token invalid")
)

type claims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

// Authenticator fetches JWKS from auth-service and validates RS256 JWT tokens.
type Authenticator struct {
	kf  keyfunc.Keyfunc
	log *slog.Logger
}

func New(ctx context.Context, authServiceURL string, refreshInterval time.Duration, log *slog.Logger) (*Authenticator, error) {
	jwksURL := strings.TrimRight(authServiceURL, "/") + "/.well-known/jwks.json"

	override := keyfunc.Override{
		RefreshInterval: refreshInterval,
	}
	kf, err := keyfunc.NewDefaultOverrideCtx(ctx, []string{jwksURL}, override)
	if err != nil {
		return nil, fmt.Errorf("init jwks: %w", err)
	}

	log.Info("jwks initialized", "url", jwksURL, "refresh_interval", refreshInterval)
	return &Authenticator{kf: kf, log: log}, nil
}

func (a *Authenticator) ValidateToken(_ context.Context, token string) (*dto.JWTClaims, error) {
	var c claims
	parsed, err := jwt.ParseWithClaims(token, &c, a.kf.Keyfunc, jwt.WithValidMethods([]string{"RS256"}))
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("%w: %s", ErrTokenInvalid, err.Error())
	}
	if !parsed.Valid {
		return nil, ErrTokenInvalid
	}
	return &dto.JWTClaims{
		Sub:  c.Subject,
		Role: c.Role,
	}, nil
}
