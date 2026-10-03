# webmemo CLI

`webmemo` is a command line client for WebMemo. It logs in to a WebMemo server
and runs an [MCP](https://modelcontextprotocol.io) server on stdio, so AI
agents such as Claude Code can search and edit your memos.

> **Server deployment required for browser login.** `webmemo login` (the
> browser flow) depends on an API change that ships with the server. Until it
> is deployed to your server, use `webmemo login --token <T>`, where `<T>` is
> the value of the `wmToken` cookie of memo.redshore.me (browser devtools,
> Application/Storage, Cookies). Token refresh and the other commands work
> without that change.

## Install

```bash
go install github.com/isutare412/web-memo/cli/cmd/webmemo@latest
```

Or, from a checkout, run `make install` in this directory (or `make install-cli` at the repo root) to install into `GOBIN` (default `~/go/bin`), or `make build` / `make build-cli` to write `cli/bin/webmemo`.

## Log in

On a machine with a browser:

```bash
webmemo login
```

This opens the browser, finishes the login through a local loopback callback,
and stores the credential in `$XDG_CONFIG_HOME/webmemo/credentials.json`
(default `~/.config/webmemo/credentials.json`). Use `--no-browser` to print the
login URL instead of opening it, and `--server <URL>` for a server other than
the default. If the browser shows an error or you cancel the Google consent,
press Ctrl-C; otherwise the CLI waits 5 minutes.

The callback is the loopback address of the machine running `webmemo login`.
Over SSH with `--no-browser`, forward the port the CLI reports
(`ssh -L <port>:127.0.0.1:<port> <host>`) so your local browser can reach it, or
skip the browser flow: run `webmemo token` on your PC and `webmemo login
--token <T>` on the server.

On a headless server, get a token on your PC and hand it over:

```bash
# on your PC
webmemo token

# on the server
webmemo login --token <T>
```

Alternatively set `WEBMEMO_TOKEN` (and `WEBMEMO_SERVER` if needed) in the
environment. The environment takes precedence over the credentials file, and
a token from the environment is never written to disk. If `WEBMEMO_SERVER`
names a different server than the saved login, the saved token is not used.

Other commands: `webmemo whoami` shows the user and token expiry, and
`webmemo logout` deletes the stored credential.

### Token refresh

Tokens are valid for 30 days. Every command and the MCP server refresh the
token automatically once it is within 7 days of expiry. To reset the expiry
on a schedule instead (for example a weekly cron job on an agent host), run:

```sh
webmemo refresh   # prints "token refreshed; expires <time>"
```

It exits non-zero when the refresh fails (no login, rejected token, network
error), so the scheduler can alert. With `WEBMEMO_TOKEN`, the new token cannot
be saved: `webmemo refresh` prints it alone on stdout and a note on stderr, so
a script can capture it and update the variable.

## MCP server

Register it with Claude Code:

```bash
claude mcp add webmemo -- webmemo mcp
```

For read-only access, expose only the read tools:

```bash
claude mcp add webmemo -- webmemo mcp --read-only
```

The server starts even when you are not logged in; tool calls then return a
message asking you to run `webmemo login`. After such a 401 (not logged in or
token expired), run `webmemo login`: the running MCP server picks up the new
token from the credentials file automatically, with no restart. If
`WEBMEMO_TOKEN` is set it overrides the file, so update it and restart the MCP
server instead. stdout carries only the MCP protocol, so any diagnostics go to
stderr.

### Tools

Read tools (also available with `--read-only`):

- `search_memos`: list memos, or run a hybrid semantic + keyword search
- `get_memo`: get one memo with its full content, tags, and version
- `list_memo_members`: list the subscribers and collaborators of a memo
- `search_tags`: search existing tag names

Write tools:

- `create_memo`: create a memo
- `update_memo`: replace the title, content, or tags of a memo
- `edit_memo`: edit a memo's content by replacing an exact string
- `delete_memo`: permanently delete a memo
- `set_memo_tags`: replace all tags of a memo
- `publish_memo`: set the publish state of a memo
- `subscribe_memo`: subscribe to or unsubscribe from another user's memo
- `request_collaboration`: request or cancel collaboration on another user's memo
- `authorize_member`: approve or reject a subscriber or collaborator of your memo
- `upload_image`: upload a local image and get a Markdown-ready URL

## Development

```bash
make test    # go test -race ./...
make check   # format, golangci-lint, tests
make generate  # regenerate the API client from the api module's OpenAPI spec
```
