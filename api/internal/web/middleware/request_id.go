package middleware

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/isutare412/web-memo/api/internal/log"
	"github.com/isutare412/web-memo/api/internal/tracing"
)

const requestIDHeader = "X-Request-Id"

func WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		span := tracing.SpanFromContext(ctx)
		spanCtx := span.SpanContext()

		// Use trace ID as request ID if possible
		id := uuid.NewString()
		if tid := spanCtx.TraceID(); tid.IsValid() {
			id = tid.String()
		}
		log.AddAttrs(ctx, slog.String("requestId", id))

		w.Header().Set(requestIDHeader, id)
		next.ServeHTTP(w, r)
	})
}
