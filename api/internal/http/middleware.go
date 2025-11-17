package http

import (
	"context"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/isutare412/web-memo/api/internal/log"
)

type contextBagKey struct{}

type contextBag struct {
	passport *passport
}

func withContextBag(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		bag := &contextBag{}
		ctx := context.WithValue(r.Context(), contextBagKey{}, bag)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func getContextBag(ctx context.Context) (*contextBag, bool) {
	bag, ok := ctx.Value(contextBagKey{}).(*contextBag)
	return bag, ok
}

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

		if bag, ok := getContextBag(r.Context()); ok {
			if bag.passport != nil {
				accessLog = accessLog.With(
					slog.String("userId", bag.passport.token.UserID.String()),
					slog.String("userType", string(bag.passport.token.UserType)),
					slog.String("userName", bag.passport.token.UserName),
					slog.String("email", bag.passport.token.Email),
				)
			}
		}

		accessLog.Log(r.Context(), log.SlogLevelAccess, "HTTP request handled")
	}

	return http.HandlerFunc(fn)
}

func recoverPanic(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("HTTP handler panicked", "recover", r, "stack", string(debug.Stack()))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
