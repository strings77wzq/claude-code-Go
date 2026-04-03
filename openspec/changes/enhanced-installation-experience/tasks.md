## Tasks

- [x] Task 1: Implement `go-code --setup` interactive wizard in Go
  - Create `cmd/go-code/setup.go` with interactive provider selection, API key input, format validation, and config file writing
  - Update `cmd/go-code/main.go` to handle `--setup` flag

- [x] Task 2: Enhance `install.sh` to call `go-code --setup` after download
  - Add setup call at end of install.sh
  - Add fallback instructions if setup fails

- [x] Task 3: Create `install.ps1` Windows PowerShell installer
  - Detect architecture (amd64)
  - Download Windows binary from GitHub Releases
  - Install to user PATH
  - Call `go-code --setup`

- [x] Task 4: Create `docs/guide/installation-for-agents.md`
  - Simplified 4-step guide for AI agents
  - Step 0: Ask user about provider
  - Step 1: Install binary
  - Step 2: Configure API key
  - Step 3: Verify and first run

- [x] Task 5: Update README.md with For LLM Agents link
  - Add link to installation-for-agents.md in Installation section

- [x] Task 6: Update `docs/guide/installation.md` with For LLM Agents chapter
  - Add section linking to installation-for-agents.md

- [x] Task 7: Build verification and docs rebuild
  - Run `go build` to verify compilation
  - Run `make docs-build` to regenerate docs site
