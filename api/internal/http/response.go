package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/isutare412/web-memo/api/internal/pkgerr"
)

type errorResponse struct {
	Msg string `json:"msg,omitempty"`
}

func responseError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		statusCode = http.StatusInternalServerError
		body       errorResponse
	)

	if kerr, ok := pkgerr.AsKnown(err); ok {
		statusCode = kerr.Code.ToHTTPStatusCode()
		if statusCode >= http.StatusInternalServerError {
			slog.Error("5xx error occurred", "error", kerr.Error(), "code", kerr.Code, "httpStatusCode", statusCode)
		}

		if kerr.Simple != nil {
			body.Msg = kerr.Simple.Error()
		}
	} else {
		slog.Error("unknown error occurred", "error", err, "httpStatusCode", statusCode)
	}

	bodyBytes, marhsalErr := json.Marshal(&body)
	if marhsalErr != nil {
		slog.Error("failed to marshap error response", "error", marhsalErr)
	}

	header := w.Header()
	header.Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(bodyBytes)
}
