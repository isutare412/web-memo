package pkgerr

import (
	"errors"
	"fmt"
	"net/http"
)

type Known struct {
	Code      Code
	Origin    error
	ClientMsg string
}

func (k Known) Error() string {
	switch {
	case k.Origin != nil:
		return k.Origin.Error()
	case k.ClientMsg != "":
		return k.ClientMsg
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
func IsErrBadRequest(err error) bool {
	if kerr, ok := AsKnown(err); ok {
		return kerr.Code == CodeBadRequest
	}
	return false
}

func IsErrNotFound(err error) bool {
	if kerr, ok := AsKnown(err); ok {
		return kerr.Code == CodeNotFound
	}
	return false
}
func IsErrAlreadyExists(err error) bool {
	if kerr, ok := AsKnown(err); ok {
		return kerr.Code == CodeAlreadyExists
	}
	return false
}

func IsErrUnauthenticated(err error) bool {
	if kerr, ok := AsKnown(err); ok {
		return kerr.Code == CodeUnauthenticated
	}
	return false
}

func IsErrPermissionDenied(err error) bool {
	if kerr, ok := AsKnown(err); ok {
		return kerr.Code == CodePermissionDenied
	}
	return false
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
