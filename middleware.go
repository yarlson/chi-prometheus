// Package chiprom provides integration of Prometheus metrics into the chi router.
// This allows for detailed monitoring of HTTP request handling within a chi-based service.
// The metrics exposed include request count, latency, and response size, categorized by status code, HTTP method, and path.
package chiprom

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

// Default bucket values for histogram metrics.
var (
	dflBuckets = []float64{300, 1200, 5000}
)

// Constant metric names used throughout the middleware.
const (
	reqsName           = "chi_requests_total"
	latencyName        = "chi_request_duration_milliseconds"
	patternReqsName    = "chi_pattern_requests_total"
	patternLatencyName = "chi_pattern_request_duration_milliseconds"
)

// Middleware encapsulates the counters and histograms for monitoring
// the number of requests, their latency, and the response size.
type Middleware struct {
	reqs    *prometheus.CounterVec
	latency *prometheus.HistogramVec
}

// NewMiddleware constructs a Middleware that records basic request metrics.
// Name parameter identifies the service, and buckets customizes latency histograms.
// It wraps the next HTTP handler, instrumenting how requests are processed.
func NewMiddleware(name string, buckets ...float64) func(next http.Handler) http.Handler {
	var m Middleware
	m.reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        reqsName,
			Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
			ConstLabels: prometheus.Labels{"service": name},
		},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(m.reqs)

	if len(buckets) == 0 {
		buckets = dflBuckets
	}
	m.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        latencyName,
		Help:        "How long it took to process the request, partitioned by status code, method and HTTP path.",
		ConstLabels: prometheus.Labels{"service": name},
		Buckets:     buckets,
	},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(m.latency)
	return m.handler
}

// A private function that defines how the basic request metrics are gathered.
func (c Middleware) handler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		c.reqs.WithLabelValues(http.StatusText(ww.Status()), r.Method, r.URL.Path).Inc()
		c.latency.WithLabelValues(http.StatusText(ww.Status()), r.Method, r.URL.Path).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
	}
	return http.HandlerFunc(fn)
}

// NewPatternMiddleware constructs a Middleware that groups requests by chi routing pattern.
// For example, patterns like /users/{firstName} can be monitored rather than specific instances like /users/bob.
// Name identifies the service and buckets customizes latency histograms.
func NewPatternMiddleware(name string, buckets ...float64) func(next http.Handler) http.Handler {
	var m Middleware
	m.reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        patternReqsName,
			Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path (with patterns).",
			ConstLabels: prometheus.Labels{"service": name},
		},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(m.reqs)

	if len(buckets) == 0 {
		buckets = dflBuckets
	}
	m.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        patternLatencyName,
		Help:        "How long it took to process the request, partitioned by status code, method and HTTP path (with patterns).",
		ConstLabels: prometheus.Labels{"service": name},
		Buckets:     buckets,
	},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(m.latency)
	return m.patternHandler
}

// A private function defining how the pattern-specific request metrics are gathered.
func (c Middleware) patternHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		rctx := chi.RouteContext(r.Context())
		routePattern := strings.Join(rctx.RoutePatterns, "")
		routePattern = strings.Replace(routePattern, "/*/", "/", -1)

		c.reqs.WithLabelValues(http.StatusText(ww.Status()), r.Method, routePattern).Inc()
		c.latency.WithLabelValues(http.StatusText(ww.Status()), r.Method, routePattern).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
	}
	return http.HandlerFunc(fn)
}
