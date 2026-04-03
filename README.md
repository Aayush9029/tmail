<p align="center">
  <img src="assets/icon.png" width="128" alt="tmail">
  <h1 align="center">tmail</h1>
  <p align="center">Disposable email in your terminal</p>
</p>

<p align="center">
  <a href="https://github.com/Aayush9029/tmail/releases/latest"><img src="https://img.shields.io/github/v/release/Aayush9029/tmail" alt="Release"></a>
  <a href="https://github.com/Aayush9029/tmail/blob/main/LICENSE"><img src="https://img.shields.io/github/license/Aayush9029/tmail" alt="License"></a>
</p>

<p align="center">
  <img src="assets/demo.gif" alt="tmail demo" width="800">
</p>

## Install

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
tmail messages          # list inbox (select to open in Safari)
tmail read 1            # read message #1
tmail read 1 --browser  # open in Safari
tmail me                # show account info
tmail delete            # delete account
```

## License

MIT
