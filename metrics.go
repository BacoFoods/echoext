package echoext

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// httpRequestsTotal counts every served request, partitioned by method,
	// templated route and status code.
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total HTTP requests served, partitioned by method, route and status code.",
	}, []string{"method", "route", "status"})

	// httpRequestDuration tracks response time distribution. Buckets cover fast
	// reads (5ms) up to slow externally-bound calls (10s).
	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request latency in seconds, partitioned by method, route and status code.",
		Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
	}, []string{"method", "route", "status"})

	// httpRequestsInFlight is the count of requests currently being processed.
	httpRequestsInFlight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "http_requests_in_flight",
		Help: "Number of HTTP requests currently being served.",
	})
)

// metricsMiddleware records Prometheus metrics for every request handled by the
// main server. OPTIONS requests (typically CORS preflights) are ignored. Routes
// are reported using their templated form (e.g. "/users/:id") to keep label
// cardinality bounded.
func metricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		// Ignore OPTIONS requests (CORS preflight, etc.).
		if req.Method == http.MethodOptions {
			return next(c)
		}

		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		start := time.Now()
		err := next(c)
		elapsed := time.Since(start).Seconds()

		// Resolve the final status code. When the handler returns an error the
		// response status may not be written yet, so fall back to the echo
		// error's code (defaulting to 500 for non-HTTP errors).
		status := c.Response().Status
		if err != nil {
			var he *echo.HTTPError
			if errors.As(err, &he) {
				status = he.Code
			} else {
				status = http.StatusInternalServerError
			}
		}

		// Use the templated route to avoid unbounded label cardinality. It is
		// empty for unmatched paths (404s).
		route := c.Path()
		if route == "" {
			route = "unmatched"
		}

		method := req.Method
		statusStr := strconv.Itoa(status)

		httpRequestsTotal.WithLabelValues(method, route, statusStr).Inc()
		httpRequestDuration.WithLabelValues(method, route, statusStr).Observe(elapsed)

		return err
	}
}

// newMetricsServer builds the dedicated HTTP server that exposes the Prometheus
// metrics endpoint on its own port, isolated from application traffic.
func newMetricsServer(c ServerConfig) *http.Server {
	mux := http.NewServeMux()
	mux.Handle(c.MetricsConfig.escapePath(), promhttp.Handler())

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", c.MetricsConfig.escapePort()),
		Handler: mux,
	}
}
