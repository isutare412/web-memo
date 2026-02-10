package middleware

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/isutare412/web-memo/api/internal/metric"
)

// ObserveMetrics records HTTP request metrics using the route pattern for low cardinality.
func ObserveMetrics(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start)

		record := GetResponseRecord(r.Context())

		statusCode := http.StatusOK
		if record != nil {
			statusCode = record.Status
		}

		// NOTE: We use route pattern instead of r.URL.Path to reduce cardinality.
		var path string
		route := mux.CurrentRoute(r)
		if route != nil {
			path, _ = route.GetPathTemplate()
		} else {
			path = r.URL.Path
		}

		metric.ObserveHTTPRequest(r.Method, path, statusCode, elapsed)
	}
	return http.HandlerFunc(fn)
}
