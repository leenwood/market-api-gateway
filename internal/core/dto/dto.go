package dto

import "net/http"

// RequestInfo carries metadata about an incoming request.
type RequestInfo struct {
	Method    string
	Path      string
	ClientIP  string
	RequestID string
}

// ResponseInfo carries metadata about an upstream response.
type ResponseInfo struct {
	StatusCode int
	Headers    http.Header
}

// JWTClaims holds the verified claims extracted from a JWT token.
type JWTClaims struct {
	Sub  string // subject — user ID
	Role string
}
