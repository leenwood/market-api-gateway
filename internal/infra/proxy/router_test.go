package proxy_test

import (
	"testing"

	"market-api-gateway/internal/core/domain"
	"market-api-gateway/internal/infra/proxy"
)

func TestRouter_GetRoute(t *testing.T) {
	r := proxy.NewRouter()

	tests := []struct {
		path        string
		method      string
		wantFound   bool
		wantService domain.ServiceName
		wantType    domain.RouteType
		wantParams  map[string]string
	}{
		// Public routes
		{"/api/v1/auth/login", "POST", true, domain.ServiceAuth, domain.Public, nil},
		{"/api/v1/products", "GET", true, domain.ServiceCatalog, domain.Public, nil},
		{"/api/v1/products/42", "GET", true, domain.ServiceCatalog, domain.Public, map[string]string{"id": "42"}},
		{"/api/v1/categories/7", "GET", true, domain.ServiceCatalog, domain.Public, map[string]string{"id": "7"}},
		{"/api/v1/analytics/sales/2024", "GET", true, domain.ServiceCatalog, domain.Public, nil},
		{"/api/v1/search", "GET", true, domain.ServiceCatalog, domain.Public, nil},
		{"/api/v1/search/autocomplete", "GET", true, domain.ServiceCatalog, domain.Public, nil},

		// Protected cart routes with path params
		{"/api/v1/cart/user-abc", "GET", true, domain.ServiceCart, domain.Protected, map[string]string{"userID": "user-abc"}},
		{"/api/v1/cart/user-abc/items", "POST", true, domain.ServiceCart, domain.Protected, map[string]string{"userID": "user-abc"}},
		{"/api/v1/cart/user-abc/items/prod-99", "DELETE", true, domain.ServiceCart, domain.Protected, map[string]string{"userID": "user-abc", "productID": "prod-99"}},
		{"/api/v1/cart/user-abc/items/prod-99", "PATCH", true, domain.ServiceCart, domain.Protected, map[string]string{"userID": "user-abc", "productID": "prod-99"}},

		// Protected order routes
		{"/api/v1/orders", "GET", true, domain.ServiceOrder, domain.Protected, nil},
		{"/api/v1/orders", "POST", true, domain.ServiceOrder, domain.Protected, nil},
		{"/api/v1/orders/order-123", "GET", true, domain.ServiceOrder, domain.Protected, map[string]string{"id": "order-123"}},
		{"/api/v1/orders/order-123/status", "PATCH", true, domain.ServiceOrder, domain.Protected, map[string]string{"id": "order-123"}},
		{"/api/v1/orders/order-123/cancel", "POST", true, domain.ServiceOrder, domain.Protected, map[string]string{"id": "order-123"}},

		// Method mismatch → not found
		{"/api/v1/auth/login", "GET", false, "", 0, nil},

		// Unknown path → not found
		{"/api/v1/unknown", "GET", false, "", 0, nil},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			route, params, found := r.GetRoute(tt.path, tt.method)
			if found != tt.wantFound {
				t.Fatalf("found=%v, want %v", found, tt.wantFound)
			}
			if !found {
				return
			}
			if route.Service != tt.wantService {
				t.Errorf("service=%q, want %q", route.Service, tt.wantService)
			}
			if route.Type != tt.wantType {
				t.Errorf("type=%v, want %v", route.Type, tt.wantType)
			}
			for k, v := range tt.wantParams {
				if got := params[k]; got != v {
					t.Errorf("param[%q]=%q, want %q", k, got, v)
				}
			}
		})
	}
}
