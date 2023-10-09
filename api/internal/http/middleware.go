package http

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func wrapResponseWriter(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
	}
	return http.HandlerFunc(fn)
}

func logRequests(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		beforeServing := time.Now()
		next.ServeHTTP(w, r)
		elapsedTime := time.Since(beforeServing)
		ww := w.(middleware.WrapResponseWriter)

		var statusCode = http.StatusOK
		if sc := ww.Status(); sc != 0 {
			statusCode = sc
		}

		accessLog := slog.With(
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
			slog.String("addr", r.RemoteAddr),
			slog.String("proto", r.Proto),
			slog.Int64("contentLength", r.ContentLength),
			slog.String("userAgent", r.UserAgent()),
			slog.Int("status", statusCode),
			slog.Int("bodyBytes", ww.BytesWritten()),
			slog.Duration("elapsed", elapsedTime),
		)

		if ct := r.Header.Get("Content-Type"); ct != "" {
			accessLog = accessLog.With(slog.String("contentType", ct))
		}

		accessLog.Info("handle HTTP request")
	}

	return http.HandlerFunc(fn)
}

func recoverPanic(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("HTTP handler panicked", "recover", r, "stack", debug.Stack())
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
