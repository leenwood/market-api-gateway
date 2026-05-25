package service

import (
	"errors"
	"strings"

	"market-api-gateway/internal/core/domain"
	"market-api-gateway/internal/core/dto"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

// ValidateRoute returns ErrUnauthorized if a protected route is accessed without claims.
func ValidateRoute(route *domain.Route, claims *dto.JWTClaims) error {
	if route.Type == domain.Protected && claims == nil {
		return ErrUnauthorized
	}
	return nil
}

// ExtractUserID returns the user ID from claims, or empty string if claims is nil.
func ExtractUserID(claims *dto.JWTClaims) string {
	if claims == nil {
		return ""
	}
	return claims.Sub
}

// CheckAccess validates that the authenticated user may access the resource.
// For cart routes (/api/v1/cart/:userID), the path userID must match claims.Sub.
func CheckAccess(route *domain.Route, claims *dto.JWTClaims, pathParams map[string]string) error {
	if route.Type == domain.Protected && claims == nil {
		return ErrUnauthorized
	}
	if route.Service == domain.ServiceCart {
		if uid, ok := pathParams["userID"]; ok {
			if !strings.EqualFold(uid, claims.Sub) {
				return ErrForbidden
			}
		}
	}
	return nil
}
