package cors

import (
	"net/http"

	rscors "github.com/rs/cors"
)

// New returns a CORS middleware configured with the given allowed origins.
func New(allowedOrigins []string) func(http.Handler) http.Handler {
	c := rscors.New(rscors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: false,
	})
	return c.Handler
}
