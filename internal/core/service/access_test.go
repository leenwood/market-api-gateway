package service_test

import (
	"testing"

	"market-api-gateway/internal/core/domain"
	"market-api-gateway/internal/core/dto"
	"market-api-gateway/internal/core/service"
)

func TestValidateRoute(t *testing.T) {
	protected := &domain.Route{Type: domain.Protected}
	public := &domain.Route{Type: domain.Public}
	claims := &dto.JWTClaims{Sub: "u1", Role: "user"}

	tests := []struct {
		name    string
		route   *domain.Route
		claims  *dto.JWTClaims
		wantErr error
	}{
		{"public no claims", public, nil, nil},
		{"public with claims", public, claims, nil},
		{"protected with claims", protected, claims, nil},
		{"protected no claims", protected, nil, service.ErrUnauthorized},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateRoute(tt.route, tt.claims)
			if err != tt.wantErr {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestExtractUserID(t *testing.T) {
	tests := []struct {
		name   string
		claims *dto.JWTClaims
		want   string
	}{
		{"nil claims", nil, ""},
		{"with sub", &dto.JWTClaims{Sub: "user-123"}, "user-123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := service.ExtractUserID(tt.claims); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCheckAccess(t *testing.T) {
	cartRoute := &domain.Route{Type: domain.Protected, Service: domain.ServiceCart}
	catalogRoute := &domain.Route{Type: domain.Protected, Service: domain.ServiceCatalog}
	publicRoute := &domain.Route{Type: domain.Public, Service: domain.ServiceCatalog}

	alice := &dto.JWTClaims{Sub: "alice", Role: "user"}

	tests := []struct {
		name       string
		route      *domain.Route
		claims     *dto.JWTClaims
		pathParams map[string]string
		wantErr    error
	}{
		{
			name:       "cart owner match",
			route:      cartRoute,
			claims:     alice,
			pathParams: map[string]string{"userID": "alice"},
			wantErr:    nil,
		},
		{
			name:       "cart owner mismatch → 403",
			route:      cartRoute,
			claims:     alice,
			pathParams: map[string]string{"userID": "bob"},
			wantErr:    service.ErrForbidden,
		},
		{
			name:       "cart no userID param",
			route:      cartRoute,
			claims:     alice,
			pathParams: map[string]string{},
			wantErr:    nil,
		},
		{
			name:       "protected catalog with claims",
			route:      catalogRoute,
			claims:     alice,
			pathParams: nil,
			wantErr:    nil,
		},
		{
			name:       "protected no claims → 401",
			route:      catalogRoute,
			claims:     nil,
			pathParams: nil,
			wantErr:    service.ErrUnauthorized,
		},
		{
			name:       "public no claims",
			route:      publicRoute,
			claims:     nil,
			pathParams: nil,
			wantErr:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CheckAccess(tt.route, tt.claims, tt.pathParams)
			if err != tt.wantErr {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}
