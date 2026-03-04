# tmail

Disposable email in your terminal. Powered by [mail.tm](https://mail.tm).

## Installation

```bash
brew install aayush9029/tap/tmail
```

Or tap first:

```bash
brew tap aayush9029/tap
brew install tmail
```

## Usage

```bash
tmail generate          # create a new disposable email (copied to clipboard)
tmail messages          # list inbox
tmail read 1            # read message #1
tmail read 1 --plain    # read as plain text
tmail read 1 --browser  # open in browser
tmail watch             # real-time notifications
tmail me                # show account info
tmail delete            # delete account
tmail domains           # list available domains
```

## Options

| Flag | Alias | Description |
|------|-------|-------------|
| `--help` | `-h` | Show help |
| `--version` | `-v` | Show version |
| `--plain` | `-p` | Plain text output (read) |
| `--browser` | `-b` | Open in browser (read) |

## How it works

1. Fetches available domains from mail.tm API
2. Generates a random address and password
3. Creates account and stores credentials in `~/.config/tmail/`
4. All inbox operations use the stored JWT token
5. Watch mode uses Server-Sent Events for real-time updates

## Requirements

- macOS or Linux

## License

MIT
