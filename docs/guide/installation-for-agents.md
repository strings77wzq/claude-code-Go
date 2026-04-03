---
title: Installation Guide for AI Agents
description: Step-by-step instructions for AI agents to help users install and configure go-code
---

# Installation Guide for AI Agents

If you're an AI Agent helping a user set up go-code, follow these steps.

## Step 0: Ask the user about their API setup

Ask the user which LLM provider they want to use:

1. **Which provider?**
   - Anthropic (Claude) → needs `ANTHROPIC_API_KEY` starting with `sk-ant-`
   - OpenAI (GPT) → needs `OPENAI_API_KEY` starting with `sk-`
   - Custom (OpenAI-compatible) → needs base URL + API key

2. **Do you already have an API key?**
   - If yes → ask them to paste it
   - If no → guide them to get one from the provider's website

## Step 1: Install go-code

Use the one-command install script:

**Linux / macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/strings77wzq/claude-code-Go/main/install.sh | bash
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/strings77wzq/claude-code-Go/main/install.ps1 | iex
```

The script will:
1. Detect OS and architecture
2. Download the appropriate pre-built binary
3. Install it to PATH
4. Automatically start the setup wizard

### Alternative: go install

```bash
go install github.com/strings77wzq/claude-code-Go/cmd/go-code@latest
```

After `go install`, run `go-code --setup` to configure.

## Step 2: Configure API key

The setup wizard (`go-code --setup`) will guide the user through:
1. Selecting their provider
2. Entering their API key (format-validated)
3. Choosing a model
4. Writing `~/.go-code/settings.json`

If the user skipped the wizard, configure manually:

```bash
mkdir -p ~/.go-code
cat > ~/.go-code/settings.json << 'EOF'
{
  "apiKey": "sk-ant-your-key-here",
  "provider": "anthropic",
  "model": "claude-sonnet-4-20250514"
}
EOF
```

## Step 3: Verify and first run

```bash
# Verify installation
go-code --help

# Run the setup wizard if not done yet
go-code --setup

# Start using go-code
go-code
```

The user should see the REPL welcome screen and be able to interact with the agent.
