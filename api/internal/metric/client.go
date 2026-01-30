package metric

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type client struct {
	gatherer prometheus.Gatherer

	httpRequestsTotal          *prometheus.CounterVec
	httpRequestDurationSeconds *prometheus.HistogramVec
}

func Init() {
	c := &client{
		gatherer:                   prometheus.DefaultGatherer,
		httpRequestsTotal:          newHTTPRequestsTotal(),
		httpRequestDurationSeconds: newHTTPRequestDurationSeconds(),
	}

	prometheus.MustRegister(c.httpRequestsTotal)
	prometheus.MustRegister(c.httpRequestDurationSeconds)

	globalObserver = c
}

func (c *client) observeHTTPRequest(method, path string, status int, duration time.Duration) {
	statusStr := strconv.Itoa(status)

	c.httpRequestsTotal.WithLabelValues(method, path, statusStr).Inc()
	c.httpRequestDurationSeconds.WithLabelValues(method, path, statusStr).Observe(duration.Seconds())
}
