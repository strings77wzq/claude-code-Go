# Privacy Policy

claude-code-Go respects your privacy. This policy explains what data we collect and how we use it.

## Data Collection

### Anonymous Telemetry (Opt-In)

If you choose to enable telemetry, we collect:

- **Feature Usage**: Which tools are used and how often
- **Error Counts**: Number of errors (without stack traces or messages)
- **Performance Metrics**: Response times and throughput
- **Version Info**: claude-code-Go version and Go version

We explicitly DO NOT collect:

- Your source code or file contents
- API keys or authentication tokens
- Session content, prompts, or responses
- Personal identifying information
- IP addresses or location data

### Local Data

The following data is stored locally on your machine:

- **Session History**: Saved in `~/.go-code/sessions/` for resume functionality
- **Configuration**: Settings in `~/.go-code/settings.json`
- **Cache**: Temporary files in `~/.go-code/cache/`

This data never leaves your machine unless you explicitly share it.

## How to Control Telemetry

### Check Status

```bash
go-code telemetry status
```

### Enable Telemetry

```bash
go-code telemetry enable
```

### Disable Telemetry

```bash
go-code telemetry disable
```

### View Collected Data

```bash
go-code telemetry view
```

## Data Storage

Telemetry data is stored locally in `~/.go-code/telemetry.jsonl`.

No data is sent to external servers unless you explicitly configure a telemetry endpoint.

## Open Source

This project is open source. You can audit the telemetry implementation:

- [internal/telemetry/client.go](https://github.com/strings77wzq/claude-code-Go/blob/main/internal/telemetry/client.go)
- [internal/telemetry/consent.go](https://github.com/strings77wzq/claude-code-Go/blob/main/internal/telemetry/consent.go)

## Contact

For privacy questions, open a GitHub Discussion or email the maintainers.

## Changes

We may update this policy. Changes will be posted on GitHub.

Last updated: April 2026
