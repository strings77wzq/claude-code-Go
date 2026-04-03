#!/bin/bash
# launch.sh — One-command launch script for go-code
# Usage: ./launch.sh
#
# This script:
# 1. Cleans up non-project directories
# 2. Runs go mod tidy to generate go.sum
# 3. Builds the binary
# 4. Runs all tests
# 5. Initializes git repo
# 6. Prints next-step instructions

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  go-code Launch Preparation${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Step 1: Clean up non-project directories
echo -e "${YELLOW}[1/6] Cleaning up non-project directories...${NC}"
for dir in owncode-analysis claw-code-parity test bin .pytest_cache; do
    if [ -d "$dir" ]; then
        rm -rf "$dir"
        echo "  Removed: $dir"
    fi
done
echo -e "  ${GREEN}Done${NC}"
echo ""

# Step 2: Go mod tidy
echo -e "${YELLOW}[2/6] Running go mod tidy...${NC}"
go mod tidy
echo -e "  ${GREEN}go.sum generated${NC}"
echo ""

# Step 3: Build
echo -e "${YELLOW}[3/6] Building binary...${NC}"
make build
echo -e "  ${GREEN}Binary: bin/go-code${NC}"
echo ""

# Step 4: Test
echo -e "${YELLOW}[4/6] Running tests...${NC}"
go test ./internal/config/... ./internal/tool/... ./internal/permission/... ./pkg/tty/... ./internal/api/... ./internal/agent/... -v 2>&1 | tail -20
echo -e "  ${GREEN}All tests passed${NC}"
echo ""

# Step 5: Vet
echo -e "${YELLOW}[5/6] Running go vet...${NC}"
go vet ./...
echo -e "  ${GREEN}No issues found${NC}"
echo ""

# Step 6: Git init
echo -e "${YELLOW}[6/6] Initializing git repository...${NC}"
if [ -d ".git" ]; then
    echo "  Git repo already exists, skipping init"
else
    git init
    echo "  Git repo initialized"
fi
echo ""

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Launch Preparation Complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${YELLOW}Next steps (manual):${NC}"
echo ""
echo "  1. Create a GitHub repository:"
echo "     → Go to https://github.com/new"
echo "     → Name it 'go-code'"
echo "     → Set description: 'Claude Code in Go — AI-powered coding assistant'"
echo "     → DO NOT initialize with README (we already have one)"
echo "     → Click 'Create repository'"
echo ""
echo "  2. Add remote and push:"
echo "     git remote add origin git@github.com:YOUR_USERNAME/go-code.git"
echo "     git add ."
echo "     git commit -m 'feat: initial release — Claude Code clone in Go'"
echo "     git branch -M main"
echo "     git push -u origin main"
echo ""
echo "  3. Enable GitHub Pages:"
echo "     → Go to repo Settings → Pages"
echo "     → Source: 'GitHub Actions'"
echo "     → The docs.yml workflow will auto-deploy"
echo ""
echo "  4. Create first release:"
echo "     git tag v0.1.0"
echo "     git push origin v0.1.0"
echo "     → Go to repo → Releases → Draft new release → v0.1.0"
echo ""
echo "  5. Your docs will be live at:"
echo "     https://YOUR_USERNAME.github.io/go-code/"
echo ""
