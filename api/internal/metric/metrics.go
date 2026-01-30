package metric

import "github.com/prometheus/client_golang/prometheus"

func newHTTPRequestsTotal() *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "webmemo",
			Subsystem: "api",
			Name:      "http_requests_total",
			Help:      "Total number of HTTP requests handled.",
		},
		[]string{"method", "path", "status"},
	)
}

func newHTTPRequestDurationSeconds() *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "webmemo",
			Subsystem: "api",
			Name:      "http_request_duration_seconds",
			Help:      "Duration of HTTP requests in seconds.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)
}
