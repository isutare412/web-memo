package auth

import (
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/isutare412/web-memo/api/internal/pkgerr"
)

const (
	queryCLICallback = "cliCallback"
	queryCLIState    = "cliState"

	cliCallbackHost = "127.0.0.1"
	cliCallbackPath = "/callback"
)

var cliStatePattern = regexp.MustCompile(`^[A-Za-z0-9_-]{16,128}$`)

// cliLogin is stored as the Redis OAuth state value when the sign-in was
// started by the webmemo CLI. It never goes into the unsigned state parameter
// sent to Google.
type cliLogin struct {
	Callback string `json:"cliCallback"`
	State    string `json:"cliState"`
}

// parseCLILogin returns nil without an error if neither CLI parameter exists.
func parseCLILogin(q url.Values) (*cliLogin, error) {
	callback, state := q.Get(queryCLICallback), q.Get(queryCLIState)
	if callback == "" && state == "" {
		return nil, nil
	}

	if err := validateCLICallback(callback); err != nil {
		return nil, err
	}

	if !cliStatePattern.MatchString(state) {
		return nil, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			ClientMsg: "cliState must be 16 to 128 characters of letters, digits, '-' or '_'",
		}
	}

	return &cliLogin{Callback: callback, State: state}, nil
}

// validateCLICallback accepts only http://127.0.0.1:<port>/callback without
// userinfo, query or fragment.
func validateCLICallback(raw string) error {
	invalid := func(err error) error {
		return pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			ClientMsg: "cliCallback must be http://127.0.0.1:<port>/callback",
			Origin:    err,
		}
	}

	if strings.ContainsAny(raw, "?#") {
		return invalid(nil)
	}

	u, err := url.Parse(raw)
	if err != nil {
		return invalid(err)
	}

	port, err := strconv.ParseUint(u.Port(), 10, 16)
	if err != nil || port == 0 {
		return invalid(err)
	}

	if u.Scheme != "http" ||
		u.User != nil ||
		u.Host != net.JoinHostPort(cliCallbackHost, u.Port()) ||
		u.Path != cliCallbackPath {
		return invalid(nil)
	}

	return nil
}
