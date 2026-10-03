package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/isutare412/web-memo/cli/internal/apiclient/gen"
	"github.com/isutare412/web-memo/cli/internal/credential"
	"github.com/isutare412/web-memo/cli/internal/login"
	"github.com/isutare412/web-memo/cli/internal/session"
)

const httpTimeout = 60 * time.Second

// browserOpener opens a URL in the user's browser. Tests replace it.
var browserOpener = openBrowser

func newFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	return fs
}

// parseFlags parses args and rejects positional arguments. It returns
// flag.ErrHelp for -h, after which the caller prints usage.
func parseFlags(fs *flag.FlagSet, args []string) error {
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() > 0 {
		return fmt.Errorf("%s: unexpected argument %q", fs.Name(), fs.Arg(0))
	}
	return nil
}

func printUsage(fs *flag.FlagSet, out io.Writer) {
	_, _ = fmt.Fprintf(out, "usage of %s:\n", fs.Name())
	fs.SetOutput(out)
	fs.PrintDefaults()
}

func runLogin(ctx context.Context, args []string, env func(string) string, store *credential.Store, out io.Writer) error {
	fs := newFlagSet("login")
	serverFlag := fs.String("server", "", "webmemo server URL (default: $WEBMEMO_SERVER, stored value, or "+credential.DefaultServer+")")
	tokenFlag := fs.String("token", "", "use this token instead of the browser login (headless)")
	noBrowser := fs.Bool("no-browser", false, "print the login URL without opening a browser")
	if err := parseFlags(fs, args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			printUsage(fs, out)
			return nil
		}
		return err
	}

	server := *serverFlag
	if server == "" {
		cred, _, err := credential.Resolve(store, env)
		if err != nil {
			return err
		}
		server = cred.Server
	}
	server, err := normalizeServer(server)
	if err != nil {
		return err
	}

	token := *tokenFlag
	if token == "" {
		lb := &login.Loopback{Server: server}
		if !*noBrowser {
			lb.OpenBrowser = browserOpener
		}
		token, err = lb.Run(ctx, out)
		if err != nil {
			return err
		}
	}

	// Verify before saving. SourceEnv means the session never writes to the
	// store; a refresh during the check is picked up through sess.Token().
	cred := credential.Credential{Server: server, Token: token}
	sess := session.New(cred, credential.SourceEnv, nil, &http.Client{Timeout: httpTimeout})
	user, err := currentUser(ctx, sess)
	if err != nil {
		return fmt.Errorf("verify token: %w", err)
	}
	cred.Token = sess.Token()
	if err := store.Save(cred); err != nil {
		return err
	}
	_, _ = fmt.Fprintf(out, "Logged in as %s <%s> on %s\n", user.UserName, user.Email, server)
	return nil
}

func runLogout(args []string, env func(string) string, store *credential.Store, out io.Writer) error {
	fs := newFlagSet("logout")
	if err := parseFlags(fs, args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			printUsage(fs, out)
			return nil
		}
		return err
	}
	if err := store.Delete(); err != nil {
		return err
	}
	_, _ = fmt.Fprintln(out, "Logged out.")
	if env("WEBMEMO_TOKEN") != "" {
		_, _ = fmt.Fprintln(out, "Note: WEBMEMO_TOKEN is still set and will keep being used.")
	}
	return nil
}

func runWhoami(ctx context.Context, args []string, env func(string) string, store *credential.Store, out io.Writer) error {
	fs := newFlagSet("whoami")
	if err := parseFlags(fs, args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			printUsage(fs, out)
			return nil
		}
		return err
	}
	cred, src, err := resolveLoggedIn(store, env)
	if err != nil {
		return err
	}
	sess := session.New(cred, src, store, &http.Client{Timeout: httpTimeout})
	user, err := currentUser(ctx, sess)
	if err != nil {
		return err
	}

	expiry := "unknown"
	if exp, ok := credential.TokenExpiry(sess.Token()); ok {
		expiry = exp.Format(time.RFC3339)
	}
	origin := "credentials file"
	if src == credential.SourceEnv {
		origin = "WEBMEMO_TOKEN"
	}
	_, _ = fmt.Fprintf(out, "Server:  %s\nUser:    %s <%s>\nExpires: %s\nSource:  %s\n",
		cred.Server, user.UserName, user.Email, expiry, origin)
	return nil
}

func runToken(args []string, env func(string) string, store *credential.Store, out io.Writer) error {
	fs := newFlagSet("token")
	if err := parseFlags(fs, args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			printUsage(fs, out)
			return nil
		}
		return err
	}
	cred, _, err := resolveLoggedIn(store, env)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintln(out, cred.Token)
	return nil
}

func resolveLoggedIn(store *credential.Store, env func(string) string) (credential.Credential, credential.Source, error) {
	cred, src, err := credential.Resolve(store, env)
	if err != nil {
		return credential.Credential{}, credential.SourceNone, err
	}
	if src == credential.SourceNone {
		return credential.Credential{}, credential.SourceNone, session.ErrUnauthorized
	}
	return cred, src, nil
}

func currentUser(ctx context.Context, sess *session.Session) (*gen.User, error) {
	res, err := sess.API().GetCurrentUserWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("request current user: %w", err)
	}
	if err := session.CheckStatus(res.HTTPResponse.StatusCode, res.Body); err != nil {
		return nil, err
	}
	if res.JSON200 == nil {
		return nil, errors.New("unexpected response from current user endpoint")
	}
	return res.JSON200, nil
}

func normalizeServer(raw string) (string, error) {
	raw = strings.TrimRight(strings.TrimSpace(raw), "/")
	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return "", fmt.Errorf("invalid server URL %q: want http(s)://host", raw)
	}
	return raw, nil
}

func openBrowser(target string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", target)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", target)
	default:
		cmd = exec.Command("xdg-open", target)
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	go func() { _ = cmd.Wait() }()
	return nil
}
