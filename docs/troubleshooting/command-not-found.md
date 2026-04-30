---
title: Command Not Found
description: Fix shell errors when the go-code binary is not on PATH after installation
---

# Command Not Found

If your shell reports `go-code: command not found`, the binary is installed but not on your `PATH`.

## Where the Binary Is Installed

### After `go install`

```bash
go install github.com/strings77wzq/claude-code-Go/cmd/go-code@latest
```

Go installs to `$GOPATH/bin`, which defaults to `$HOME/go/bin`:

```bash
go env GOPATH
# Typical: /home/yourname/go
```

The binary is at `$HOME/go/bin/go-code`.

### After `go build`

```bash
git clone https://github.com/strings77wzq/claude-code-Go.git
cd claude-code-Go
go build -o go-code ./cmd/go-code
```

The binary is `./go-code` in the current directory. Move it to your PATH:

```bash
mv go-code ~/go/bin/
```

### Using Makefile

```bash
make build    # builds to bin/go-code
make install  # installs to $GOPATH/bin/go-code
```

## Fixing PATH

### Bash (`~/.bashrc`)

```bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

### Zsh (`~/.zshrc`)

```bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### Fish (`~/.config/fish/config.fish`)

```fish
echo 'set -gx PATH $HOME/go/bin $PATH' >> ~/.config/fish/config.fish
source ~/.config/fish/config.fish
```

### Verify

```bash
which go-code
# Expected: /home/yourname/go/bin/go-code

go-code doctor
# Should print system readiness info
```

## Alternative Locations

System-wide: `sudo cp go-code /usr/local/bin/`
User-local: `mkdir -p ~/.local/bin && cp go-code ~/.local/bin/`

## Shell Cache

If PATH is correct but the command is still not found:
`hash -r` (Bash) or `rehash` (Zsh), or start a new terminal.

## Related

- [Installation Guide](../guide/installation.md) — Full install instructions
- [Common Issues](common-issues.md) — Other installation and config problems
