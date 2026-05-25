package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"market-api-gateway/internal/core/port"
)

// HealthHandler responds to liveness probes.
func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// ReadyHandler pings all registered upstreams; returns 503 if any fail.
func ReadyHandler(pingers []port.Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		type result struct {
			name string
			err  error
		}
		ch := make(chan result, len(pingers))
		for _, p := range pingers {
			p := p
			go func() {
				ch <- result{name: p.Name(), err: p.Ping(ctx)}
			}()
		}

		failed := make(map[string]string)
		for range pingers {
			res := <-ch
			if res.err != nil {
				failed[res.name] = res.err.Error()
			}
		}

		if len(failed) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]any{"status": "unavailable", "details": failed})
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
