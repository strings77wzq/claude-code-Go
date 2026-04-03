# Installation

This guide covers how to install go-code from source or pre-built binaries.

## Prerequisites

- Go 1.23 or later
- Python 3.x (optional, only needed for the test harness)

## From Source

### Clone the Repository

```bash
git clone https://github.com/user/go-code.git
cd go-code
```

### Build

```bash
make build
```

This creates the binary at `bin/go-code`.

### Install to PATH

```bash
go install ./cmd/go-code
```

This installs the binary to `$GOPATH/bin` (typically `~/go/bin`).

## Pre-built Binaries

Download the appropriate binary for your platform from the releases page:

| Platform | Architecture | Filename |
|----------|--------------|-----------|
| Linux | amd64 | go-code-linux-amd64 |
| macOS | amd64 | go-code-darwin-amd64 |
| macOS | arm64 | go-code-darwin-arm64 |

```bash
# Example: Download and install Linux amd64
curl -L -o go-code https://github.com/user/go-code/releases/latest/download/go-code-linux-amd64
chmod +x go-code
sudo mv go-code /usr/local/bin/
```

## Verify Installation

```bash
# Check the binary exists
ls -la bin/go-code

# Run with --help
./bin/go-code --help
```

## Next Steps

- [Quick Start](quick-start.md) - Run your first command
- [Configuration](config.md) - Set up API keys and preferences