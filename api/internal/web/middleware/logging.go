package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// AccessLog logs HTTP request details after the response is written.
func AccessLog(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		beforeServing := time.Now()
		next.ServeHTTP(w, r)
		elapsedTime := time.Since(beforeServing)

		record := GetResponseRecord(r.Context())

		statusCode := http.StatusOK
		bytesWritten := 0
		if record != nil {
			statusCode = record.Status
			bytesWritten = record.BytesWritten
		}

		accessLog := slog.With(
			slog.String("logType", "accessLog"),
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
			slog.String("addr", r.RemoteAddr),
			slog.String("proto", r.Proto),
			slog.Int64("contentLength", r.ContentLength),
			slog.String("userAgent", r.UserAgent()),
			slog.Int("status", statusCode),
			slog.Int("bodyBytes", bytesWritten),
			slog.Duration("elapsed", elapsedTime),
		)

		if ct := r.Header.Get("Content-Type"); ct != "" {
			accessLog = accessLog.With(slog.String("contentType", ct))
		}

		if bag, ok := GetContextBag(r.Context()); ok {
			if bag.Passport != nil {
				accessLog = accessLog.With(
					slog.String("userId", bag.Passport.Token.UserID.String()),
					slog.String("userType", string(bag.Passport.Token.UserType)),
					slog.String("userName", bag.Passport.Token.UserName),
					slog.String("email", bag.Passport.Token.Email),
				)
			}
		}

		accessLog.InfoContext(r.Context(), "HTTP request handled")
	}
	return http.HandlerFunc(fn)
}
