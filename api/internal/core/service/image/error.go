package image

import (
	"fmt"

	"github.com/isutare412/imageer/pkg/gateway"

	"github.com/isutare412/web-memo/api/internal/pkgerr"
)

func imageerError(statusCode int, errResp *gateway.ErrorResponse) error {
	var codeName, message string
	if errResp != nil {
		codeName = errResp.CodeName
		message = errResp.Message
	}

	return pkgerr.Known{
		Code:      pkgerr.CodeFromHTTPStatusCode(statusCode),
		Origin:    fmt.Errorf("imageer error: [%s] %s (status %d)", codeName, message, statusCode),
		ClientMsg: message,
	}
}
