package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"

	"github.com/isutare412/web-memo/api/internal/core/port"
)

type healthHandler struct {
	pingers []port.Pinger
}

func newHealthHandler(pingers []port.Pinger) *healthHandler {
	return &healthHandler{
		pingers: pingers,
	}
}

func (h *healthHandler) router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/livez", h.checkLiveness)
	r.Get("/readyz", h.checkReadiness)

	return r
}

func (h *healthHandler) checkLiveness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *healthHandler) checkReadiness(w http.ResponseWriter, r *http.Request) {
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
