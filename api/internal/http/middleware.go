package http

import (
	"context"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/isutare412/web-memo/api/internal/log"
	"github.com/isutare412/web-memo/api/internal/metric"
	"github.com/isutare412/web-memo/api/internal/trace"
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

		statusCode := http.StatusOK
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

		accessLog.InfoContext(r.Context(), "HTTP request handled")
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

func observeMetrics(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start)

		ww := w.(middleware.WrapResponseWriter)
		statusCode := ww.Status()
		if statusCode == 0 {
			statusCode = http.StatusOK
		}

		// NOTE: We use route pattern instead of r.URL.Path to reduce cardinality
		path := chi.RouteContext(ctx).RoutePattern()
		metric.ObserveHTTPRequest(r.Method, path, statusCode, elapsed)
	}
	return http.HandlerFunc(fn)
}

func withLogAttrContext(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = log.WithAttrContext(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func withTrace(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = trace.ExtractFromHTTPHeader(ctx, r.Header)

		ctx, span := trace.StartSpan(ctx, "http.middleware.withTrace")
		defer span.End()

		// NOTE: If sampling decision is "not sampled", trace id will be zero-value.
		spanCtx := span.SpanContext()
		if traceID := spanCtx.TraceID().String(); traceID != "" {
			log.AddAttrs(ctx, slog.String("traceId", traceID))
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
