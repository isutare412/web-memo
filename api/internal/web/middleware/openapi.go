package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gorilla/mux"
	oapimiddleware "github.com/oapi-codegen/nethttp-middleware"

	"github.com/isutare412/web-memo/api/internal/pkgerr"
	"github.com/isutare412/web-memo/api/internal/web/gen"
)

// WithOpenAPIValidator returns a middleware that validates incoming requests against
// the embedded OpenAPI specification.
func WithOpenAPIValidator() mux.MiddlewareFunc {
	swagger, err := gen.GetSwagger()
	if err != nil {
		panic(fmt.Errorf("getting swagger spec: %w", err))
	}

	// Disable schema error details in responses.
	openapi3.SchemaErrorDetailsDisabled = true

	return oapimiddleware.OapiRequestValidatorWithOptions(swagger, &oapimiddleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
		},
		ErrorHandlerWithOpts:  handleOpenAPIError,
		SilenceServersWarning: true,
		DoNotValidateServers:  true,
	})
}

func handleOpenAPIError(ctx context.Context, err error, w http.ResponseWriter, r *http.Request,
	opts oapimiddleware.ErrorHandlerOpts,
) {
	var (
		summary = "Failed to validate request"
		detail  = err.Error()
	)

	var reqErr *openapi3filter.RequestError
	if errors.As(err, &reqErr) {
		summary = reqErr.Reason
		if reqErr.Err != nil {
			detail = reqErr.Err.Error()
		}
	}

	gen.RespondError(w, r, pkgerr.Known{
		Code:      pkgerr.CodeFromHTTPStatusCode(opts.StatusCode),
		ClientMsg: fmt.Sprintf("%s: %s", summary, detail),
	})
}
