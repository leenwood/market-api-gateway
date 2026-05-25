package adapter

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// HTTPPinger implements port.Pinger by issuing a GET /health to an upstream.
type HTTPPinger struct {
	name   string
	url    string
	client *http.Client
}

func NewHTTPPinger(name, baseURL string, client *http.Client) *HTTPPinger {
	return &HTTPPinger{
		name:   name,
		url:    strings.TrimRight(baseURL, "/") + "/health",
		client: client,
	}
}

func (p *HTTPPinger) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.url, nil)
	if err != nil {
		return fmt.Errorf("build ping request: %w", err)
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("ping %s: %w", p.name, err)
	}
	resp.Body.Close()
	if resp.StatusCode >= 500 {
		return fmt.Errorf("upstream %s returned %d", p.name, resp.StatusCode)
	}
	return nil
}

func (p *HTTPPinger) Name() string { return p.name }
