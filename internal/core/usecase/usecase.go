package usecase

import (
	"context"
	"fmt"

	"market-api-gateway/internal/core/domain"
	"market-api-gateway/internal/core/dto"
	"market-api-gateway/internal/core/port"
	"market-api-gateway/internal/core/service"
)

// ValidateJWT validates a Bearer token and returns claims.
type ValidateJWT struct {
	auth port.Authenticator
}

func NewValidateJWT(auth port.Authenticator) *ValidateJWT {
	return &ValidateJWT{auth: auth}
}

func (uc *ValidateJWT) Execute(ctx context.Context, token string) (*dto.JWTClaims, error) {
	claims, err := uc.auth.ValidateToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("jwt validation: %w", err)
	}
	return claims, nil
}

// CheckAccess validates that a user may access a resource.
type CheckAccess struct{}

func NewCheckAccess() *CheckAccess { return &CheckAccess{} }

func (uc *CheckAccess) Execute(route *domain.Route, claims *dto.JWTClaims, pathParams map[string]string) error {
	return service.CheckAccess(route, claims, pathParams)
}
