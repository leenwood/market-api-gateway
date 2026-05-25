package proxy

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"

	"market-api-gateway/internal/core/domain"
	"market-api-gateway/internal/infra/outbound/httpclient"
	"market-api-gateway/internal/platform/metrics"
)

// Proxy holds one httputil.ReverseProxy per upstream service.
type Proxy struct {
	proxies map[domain.ServiceName]*httputil.ReverseProxy
	log     *slog.Logger
}

type ServiceClients struct {
	AuthURL    string
	CatalogURL string
	CartURL    string
	OrderURL   string
	Clients    map[domain.ServiceName]*httpclient.Client
}

func NewProxy(svc ServiceClients, m *metrics.Metrics, log *slog.Logger) (*Proxy, error) {
	p := &Proxy{
		proxies: make(map[domain.ServiceName]*httputil.ReverseProxy, 4),
		log:     log,
	}

	targets := map[domain.ServiceName]string{
		domain.ServiceAuth:    svc.AuthURL,
		domain.ServiceCatalog: svc.CatalogURL,
		domain.ServiceCart:    svc.CartURL,
		domain.ServiceOrder:   svc.OrderURL,
	}

	for name, rawURL := range targets {
		target, err := url.Parse(rawURL)
		if err != nil {
			return nil, fmt.Errorf("parse upstream url for %s: %w", name, err)
		}
		client := svc.Clients[name]
		rp := &httputil.ReverseProxy{
			Rewrite: func(pr *httputil.ProxyRequest) {
				pr.SetURL(target)
				pr.SetXForwarded()
			},
			Transport: client,
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
				log.ErrorContext(r.Context(), "proxy error", "upstream", name, "err", err)
				http.Error(w, `{"error":"bad gateway"}`, http.StatusBadGateway)
			},
		}
		p.proxies[name] = rp
	}
	return p, nil
}

// ServeHTTP forwards the request to the appropriate upstream based on route.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request, service domain.ServiceName) {
	rp, ok := p.proxies[service]
	if !ok {
		http.Error(w, `{"error":"no upstream"}`, http.StatusBadGateway)
		return
	}
	rp.ServeHTTP(w, r)
}
