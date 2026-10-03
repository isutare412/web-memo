// Command webmemo is the web-memo command line client.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/isutare412/web-memo/cli/internal/credential"
)

const usage = `usage: webmemo <command> [flags]

commands:
  login    log in via the browser (or --token T) and store the credential
  logout   delete the stored credential
  whoami   show the logged-in user and token expiry
  token    print the current token

Run "webmemo <command> -h" for command flags.
`

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	store, err := credential.DefaultStore()
	if err != nil {
		fmt.Fprintln(os.Stderr, "webmemo:", err)
		os.Exit(1)
	}
	if err := run(ctx, os.Args[1:], os.Getenv, store, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "webmemo:", err)
		os.Exit(1)
	}
}

// run dispatches args to a subcommand.
func run(ctx context.Context, args []string, env func(string) string, store *credential.Store, out io.Writer) error {
	if len(args) == 0 {
		return errors.New("missing command\n\n" + usage)
	}
	cmd, rest := args[0], args[1:]
	switch cmd {
	case "login":
		return runLogin(ctx, rest, env, store, out)
	case "logout":
		return runLogout(rest, env, store, out)
	case "whoami":
		return runWhoami(ctx, rest, env, store, out)
	case "token":
		return runToken(rest, env, store, out)
	// "mcp" is added by the MCP server task.
	case "-h", "-help", "--help", "help":
		_, _ = io.WriteString(out, usage)
		return nil
	default:
		return fmt.Errorf("unknown command %q\n\n%s", cmd, usage)
	}
}
