package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/isutare412/web-memo/api/internal/pkgerr"
	"github.com/isutare412/web-memo/api/internal/web/gen"
)

// RecoverPanic recovers from panics in HTTP handlers and returns a 500 JSON error response.
func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			v := recover()
			if v == nil {
				return
			}

			slog.Error("HTTP handler panicked", "recover", v, "stack", string(debug.Stack()))

			gen.RespondError(w, r, pkgerr.Known{
				Origin: fmt.Errorf("panic recover: %v", v),
			})
		}()

		next.ServeHTTP(w, r)
	})
}
