package metric

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var globalObserver observer = noopObserver{}

func Gatherer() prometheus.Gatherer {
	if o, ok := globalObserver.(*client); ok {
		return o.gatherer
	}
	return nil
}

func ObserveHTTPRequest(method, path string, status int, duration time.Duration) {
	globalObserver.observeHTTPRequest(method, path, status, duration)
}

type observer interface {
	observeHTTPRequest(method, path string, status int, duration time.Duration)
}

type noopObserver struct{}

func (noopObserver) observeHTTPRequest(method, path string, status int, duration time.Duration) {}
