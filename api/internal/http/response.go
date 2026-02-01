package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/isutare412/web-memo/api/internal/pkgerr"
)

type errorResponse struct {
	Msg string `json:"msg"`
}

func responseError(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()

	var (
		statusCode = http.StatusInternalServerError
		body       errorResponse
	)

	if kerr, ok := pkgerr.AsKnown(err); ok {
		statusCode = kerr.Code.ToHTTPStatusCode()
		switch {
		case statusCode >= http.StatusInternalServerError:
			slog.ErrorContext(ctx, "5xx error occurred", "error", err, "code", kerr.Code,
				"httpStatusCode", statusCode)
		case statusCode >= http.StatusBadRequest:
			slog.InfoContext(ctx, "4xx error occurred", "error", err, "code", kerr.Code,
				"httpStatusCode", statusCode)
		}

		if kerr.ClientMsg != "" {
			body.Msg = kerr.ClientMsg
		}
	} else {
		slog.ErrorContext(ctx, "unknown error occurred", "error", err, "httpStatusCode", statusCode)
	}

	if body.Msg == "" {
		body.Msg = http.StatusText(statusCode)
	}

	bodyBytes, marhsalErr := json.Marshal(&body)
	if marhsalErr != nil {
		slog.ErrorContext(ctx, "failed to marshap error response", "error", marhsalErr)
	}

	header := w.Header()
	header.Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if _, err := w.Write(bodyBytes); err != nil {
		slog.ErrorContext(ctx, "failed to write error response as HTTP response", "error", err)
		return
	}
}

func responsePassportError(w http.ResponseWriter, r *http.Request) {
	responseError(w, r, pkgerr.Known{
		Code:      pkgerr.CodeUnauthenticated,
		ClientMsg: "need token",
	})
}

func responseJSON(w http.ResponseWriter, obj any) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		slog.Error("failed to marshal HTTP response", "error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(bytes); err != nil {
		slog.Error("failed to write JSON bytes as HTTP response", "error", err)
		return
	}
}

func responseStatusCode(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}
