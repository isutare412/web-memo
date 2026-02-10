package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"golang.org/x/sync/errgroup"
)

// Liveness returns 200 OK to indicate the server is alive.
func (h *Handler) Liveness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Readiness pings all pingers in parallel and returns 503 if any fail.
func (h *Handler) Readiness(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	eg, ctx := errgroup.WithContext(ctx)

	for _, pinger := range h.pingers {
		pinger := pinger

		eg.Go(func() error {
			if err := pinger.Ping(ctx); err != nil {
				return fmt.Errorf("pinging target(%s): %w", pinger.Name(), err)
			}
			return nil
		})
	}

	code := http.StatusOK
	if err := eg.Wait(); err != nil {
		slog.Error("ping failed", "error", err)
		code = http.StatusServiceUnavailable
	}

	w.WriteHeader(code)
}
