package pkgerr

import "fmt"

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

type Code int

const (
	CodeUnspecified Code = iota
	CodeBadRequest
	CodeNotFound
	CodeAlreadyExists
	CodeUnauthenticated
	CodePermissionDenied
)
