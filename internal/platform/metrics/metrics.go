package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	Registry *prometheus.Registry

	HTTPRequestsTotal        *prometheus.CounterVec
	HTTPRequestDuration      *prometheus.HistogramVec
	HTTPOutboundRequestTotal *prometheus.CounterVec
}

func New() *Metrics {
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	m := &Metrics{
		Registry: reg,

		HTTPRequestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "platform_http_requests_total",
			Help: "Total number of inbound HTTP requests.",
		}, []string{"method", "path", "status"}),

		HTTPRequestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "platform_http_request_duration_seconds",
			Help:    "Inbound HTTP request duration in seconds.",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
		}, []string{"method", "path"}),

		HTTPOutboundRequestTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "platform_http_outbound_requests_total",
			Help: "Total number of outbound HTTP requests.",
		}, []string{"target", "method", "status"}),
	}

	reg.MustRegister(
		m.HTTPRequestsTotal,
		m.HTTPRequestDuration,
		m.HTTPOutboundRequestTotal,
	)

	return m
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.Registry, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	})
}
