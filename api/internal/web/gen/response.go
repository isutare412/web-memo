package gen

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/isutare412/web-memo/api/internal/pkgerr"
)

func RespondJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("Failed to json encode response", "error", err)
	}
}

func RespondNoContent(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func RespondError(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()

	var (
		statusCode = http.StatusInternalServerError
		msg        = ""
	)

	var (
		kerr              pkgerr.Known
		cookieParamErr    *UnescapedCookieParamError
		unmarshalParamErr *UnmarshalingParamError
		requiredParamErr  *RequiredParamError
		invalidParamErr   *InvalidParamFormatError
		tooManyParamErr   *TooManyValuesForParamError
		requiredHeaderErr *RequiredHeaderError
	)
	switch {
	case errors.As(err, &kerr):
		statusCode = kerr.Code.ToHTTPStatusCode()
		msg = kerr.ClientMsg
	case errors.As(err, &cookieParamErr):
		statusCode = http.StatusUnprocessableEntity
		msg = cookieParamErr.Error()
	case errors.As(err, &unmarshalParamErr):
		statusCode = http.StatusUnprocessableEntity
		msg = unmarshalParamErr.Error()
	case errors.As(err, &requiredParamErr):
		statusCode = http.StatusBadRequest
		msg = requiredParamErr.Error()
	case errors.As(err, &invalidParamErr):
		statusCode = http.StatusUnprocessableEntity
		msg = invalidParamErr.Error()
	case errors.As(err, &tooManyParamErr):
		statusCode = http.StatusUnprocessableEntity
		msg = tooManyParamErr.Error()
	case errors.As(err, &requiredHeaderErr):
		statusCode = http.StatusBadRequest
		msg = requiredHeaderErr.Error()
	}

	if msg == "" {
		msg = http.StatusText(statusCode)
	}

	entry := slog.With("statusCode", statusCode, "error", err, "clientMsg", msg)
	switch {
	case statusCode >= 500:
		entry.ErrorContext(ctx, "5xx server error occurred")
	case statusCode >= 400:
		entry.WarnContext(ctx, "4xx client error occurred")
	}

	RespondJSON(w, statusCode, ErrorResponse{Msg: msg})
}
