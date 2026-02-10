package middleware

import (
	"context"
	"net/http"
)

type responseRecordKey struct{}

// ResponseRecord tracks the status code and bytes written by the response writer.
type ResponseRecord struct {
	Status       int
	BytesWritten int
}

// responseRecorder wraps http.ResponseWriter to record status code and bytes written.
type responseRecorder struct {
	http.ResponseWriter
	record *ResponseRecord
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.record.Status = code
	rr.ResponseWriter.WriteHeader(code)
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	n, err := rr.ResponseWriter.Write(b)
	rr.record.BytesWritten += n
	return n, err
}

// WithResponseRecord wraps the response writer to record status code and bytes written.
// This middleware must run before AccessLog and ObserveMetrics.
func WithResponseRecord(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		record := &ResponseRecord{Status: http.StatusOK}
		recorder := &responseRecorder{
			ResponseWriter: w,
			record:         record,
		}
		ctx := context.WithValue(r.Context(), responseRecordKey{}, record)
		next.ServeHTTP(recorder, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// GetResponseRecord retrieves the ResponseRecord from the request context.
func GetResponseRecord(ctx context.Context) *ResponseRecord {
	record, _ := ctx.Value(responseRecordKey{}).(*ResponseRecord)
	return record
}
