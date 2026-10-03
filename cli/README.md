# webmemo CLI

`webmemo` is a command line client for WebMemo. It logs in to a WebMemo server
and runs an [MCP](https://modelcontextprotocol.io) server on stdio, so AI
agents such as Claude Code can search and edit your memos.

> **Server deployment required.** Browser login, `webmemo token` on a
> logged-in machine, and token refresh depend on API changes that ship with the
> server. Until that api change is deployed to your server, only
> `webmemo login --token <T>` works.

## Install

```bash
go install github.com/isutare412/web-memo/cli/cmd/webmemo@latest
```

Or build from this directory with `make build` (writes `bin/webmemo`).

## Log in

On a machine with a browser:

```bash
webmemo login
```

This opens the browser, finishes the login through a local loopback callback,
and stores the credential in `$XDG_CONFIG_HOME/webmemo/credentials.json`
(default `~/.config/webmemo/credentials.json`). Use `--no-browser` to print the
login URL instead of opening it, and `--server <URL>` for a server other than
the default.

On a headless server, get a token on your PC and hand it over:

```bash
# on your PC
webmemo token

# on the server
webmemo login --token <T>
```

Alternatively set `WEBMEMO_TOKEN` (and `WEBMEMO_SERVER` if needed) in the
environment. The environment takes precedence over the credentials file, and
a token from the environment is never written to disk.

Other commands: `webmemo whoami` shows the user and token expiry, and
`webmemo logout` deletes the stored credential.

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
message asking you to run `webmemo login`. stdout carries only the MCP
protocol, so any diagnostics go to stderr.

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
