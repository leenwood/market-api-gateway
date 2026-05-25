package domain

// RouteType indicates whether a route requires JWT authentication.
type RouteType int

const (
	Public    RouteType = iota
	Protected RouteType = iota
)

// ServiceName is the logical name of an upstream service.
type ServiceName string

const (
	ServiceAuth    ServiceName = "auth"
	ServiceCatalog ServiceName = "catalog"
	ServiceCart    ServiceName = "cart"
	ServiceOrder   ServiceName = "order"
)

// Route describes a single API route entry in the routing table.
type Route struct {
	Method  string
	Pattern string // e.g. /api/v1/cart/:userID
	Type    RouteType
	Service ServiceName
}

// HeaderRule describes a header to inject into upstream requests.
type HeaderRule struct {
	Name  string
	Value string
}

// Routes is the complete static routing table derived from api-gateway.md.
var Routes = []Route{
	// --- Public: auth-service ---
	{Method: "POST", Pattern: "/api/v1/auth/register", Type: Public, Service: ServiceAuth},
	{Method: "POST", Pattern: "/api/v1/auth/login", Type: Public, Service: ServiceAuth},
	{Method: "POST", Pattern: "/api/v1/auth/refresh", Type: Public, Service: ServiceAuth},

	// --- Public: market-core ---
	{Method: "GET", Pattern: "/api/v1/products", Type: Public, Service: ServiceCatalog},
	{Method: "GET", Pattern: "/api/v1/products/:id", Type: Public, Service: ServiceCatalog},
	{Method: "GET", Pattern: "/api/v1/categories", Type: Public, Service: ServiceCatalog},
	{Method: "GET", Pattern: "/api/v1/categories/:id", Type: Public, Service: ServiceCatalog},
	{Method: "GET", Pattern: "/api/v1/search", Type: Public, Service: ServiceCatalog},
	{Method: "GET", Pattern: "/api/v1/search/autocomplete", Type: Public, Service: ServiceCatalog},
	{Method: "GET", Pattern: "/api/v1/analytics/", Type: Public, Service: ServiceCatalog}, // wildcard prefix

	// --- Protected: auth-service ---
	{Method: "POST", Pattern: "/api/v1/auth/logout", Type: Protected, Service: ServiceAuth},
	{Method: "GET", Pattern: "/api/v1/auth/me", Type: Protected, Service: ServiceAuth},
	{Method: "DELETE", Pattern: "/api/v1/auth/me", Type: Protected, Service: ServiceAuth},

	// --- Protected: market-core ---
	{Method: "POST", Pattern: "/api/v1/products", Type: Protected, Service: ServiceCatalog},
	{Method: "PUT", Pattern: "/api/v1/products/:id", Type: Protected, Service: ServiceCatalog},
	{Method: "DELETE", Pattern: "/api/v1/products/:id", Type: Protected, Service: ServiceCatalog},
	{Method: "POST", Pattern: "/api/v1/categories", Type: Protected, Service: ServiceCatalog},
	{Method: "DELETE", Pattern: "/api/v1/categories/:id", Type: Protected, Service: ServiceCatalog},
	{Method: "GET", Pattern: "/api/v1/favorites", Type: Protected, Service: ServiceCatalog},
	{Method: "POST", Pattern: "/api/v1/favorites", Type: Protected, Service: ServiceCatalog},
	{Method: "DELETE", Pattern: "/api/v1/favorites", Type: Protected, Service: ServiceCatalog},

	// --- Protected: marketplace-bucket ---
	{Method: "GET", Pattern: "/api/v1/cart/:userID", Type: Protected, Service: ServiceCart},
	{Method: "DELETE", Pattern: "/api/v1/cart/:userID", Type: Protected, Service: ServiceCart},
	{Method: "POST", Pattern: "/api/v1/cart/:userID/items", Type: Protected, Service: ServiceCart},
	{Method: "DELETE", Pattern: "/api/v1/cart/:userID/items/:productID", Type: Protected, Service: ServiceCart},
	{Method: "PATCH", Pattern: "/api/v1/cart/:userID/items/:productID", Type: Protected, Service: ServiceCart},

	// --- Protected: order-service ---
	{Method: "GET", Pattern: "/api/v1/orders", Type: Protected, Service: ServiceOrder},
	{Method: "POST", Pattern: "/api/v1/orders", Type: Protected, Service: ServiceOrder},
	{Method: "GET", Pattern: "/api/v1/orders/:id", Type: Protected, Service: ServiceOrder},
	{Method: "PATCH", Pattern: "/api/v1/orders/:id/status", Type: Protected, Service: ServiceOrder},
	{Method: "POST", Pattern: "/api/v1/orders/:id/cancel", Type: Protected, Service: ServiceOrder},
}
