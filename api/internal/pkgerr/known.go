package pkgerr

import (
	"errors"
	"fmt"
	"net/http"
)

type Known struct {
	Code   Code
	Simple error
	Origin error
}

func (k Known) Error() string {
	switch {
	case k.Origin != nil:
		return k.Origin.Error()
	case k.Simple != nil:
		return k.Simple.Error()
	}
	return fmt.Sprintf("errno(%d)", k.Code)
}

func (k Known) Unwrap() error {
	return k.Origin
}

func AsKnown(err error) (Known, bool) {
	var kerr Known
	switch errors.As(err, &kerr) {
	case true:
		return kerr, true
	default:
		return Known{}, false
	}
}

type Code int

const (
	CodeUnspecified Code = iota
	CodeBadRequest
	CodeNotFound
	CodeAlreadyExists
	CodeUnauthenticated
	CodePermissionDenied
)

func (c Code) ToHTTPStatusCode() int {
	var status int
	switch c {
	case CodeUnspecified:
		status = http.StatusInternalServerError
	case CodeBadRequest:
		status = http.StatusBadRequest
	case CodeNotFound:
		status = http.StatusNotFound
	case CodeAlreadyExists:
		status = http.StatusConflict
	case CodeUnauthenticated:
		status = http.StatusUnauthorized
	case CodePermissionDenied:
		status = http.StatusForbidden
	default:
		status = http.StatusInternalServerError
	}
	return status
}
