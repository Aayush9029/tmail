# tmail

`tmail` is a Go CLI that makes a disposable email address on mail.tm and reads its inbox in the terminal. It ships through the Homebrew tap `aayush9029/tap`. Status: shipped.

## Commands

| Task | Command |
|------|---------|
| Build | `go build -ldflags "-s -w -X main.version=dev" -o tmail .` |
| Run | `go run . generate` |
| Inbox | `go run . messages` |
| Release | Load the `oss-cli-release` skill |

## Architecture

- `cmd/`: subcommands `generate`, `messages`, `read`, `me`, `delete`.
- `internal/api`: client for `https://api.mail.tm` (domains, accounts, token, messages).
- `internal/config`: the current account in `~/.config/tmail/account.json`.
- `internal/ui`: colors and the Bubble Tea inbox.

## Rules

- `generate` copies the address to the clipboard. `messages` and `read --browser` open a message in Safari.
- The GoReleaser workflow (`.github/workflows/release.yml`, `.goreleaser.yml`) fails on each tag push. Update `../../homebrew-tap/Formula/tmail.rb` by hand, as for the other tools. The formula builds from the tag tarball with `-X main.version=#{version}`.
- The built `tmail` binary at the root is in `.gitignore`. Do not commit it.
