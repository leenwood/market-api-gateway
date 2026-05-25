package mapper

import (
	"net/http"

	"market-api-gateway/internal/core/dto"
)

func MapRequestToDTO(r *http.Request, requestID string) *dto.RequestInfo {
	ip := r.RemoteAddr
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ip = xff
	}
	return &dto.RequestInfo{
		Method:    r.Method,
		Path:      r.URL.Path,
		ClientIP:  ip,
		RequestID: requestID,
	}
}

func MapResponseToDTO(statusCode int, headers http.Header) *dto.ResponseInfo {
	return &dto.ResponseInfo{
		StatusCode: statusCode,
		Headers:    headers,
	}
}
