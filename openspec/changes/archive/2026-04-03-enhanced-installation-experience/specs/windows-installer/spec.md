## ADDED Requirements

### Requirement: Windows PowerShell Installer
A PowerShell script `install.ps1` SHALL provide one-command installation for Windows users.

#### Scenario: Windows user runs install script
- **WHEN** user runs `irm <url> | iex` in PowerShell
- **THEN** the script downloads the Windows binary and installs it to a PATH directory

#### Scenario: Binary is installed
- **WHEN** the script completes
- **THEN** `go-code.exe` is available in the user's PATH and `go-code --setup` is called
