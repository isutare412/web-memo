package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

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

	var isPingFail bool
	for _, pinger := range h.pingers {
		if err := pinger.Ping(ctx); err != nil {
			slog.Error("ping failed", "target", pinger.Name())
			isPingFail = true
			break
		}
	}

	code := http.StatusOK
	if isPingFail {
		code = http.StatusServiceUnavailable
	}

	w.WriteHeader(code)
}
