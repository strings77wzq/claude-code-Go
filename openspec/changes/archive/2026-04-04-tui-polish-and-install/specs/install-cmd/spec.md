## MODIFIED Requirements

### Requirement: go install support
The project SHALL support `go install ./cmd/go-code` for one-command installation.

#### Scenario: Install command
- **WHEN** a user runs `go install ./cmd/go-code`
- **THEN** the binary is installed to $GOPATH/bin and can be run as `go-code`
