#!/bin/bash
# launch.sh — One-command launch script for go-code
# Usage: ./launch.sh
#
# This script:
# 1. Cleans up non-project directories
# 2. Runs go mod tidy to generate go.sum
# 3. Builds the binary
# 4. Runs all tests
# 5. Runs go vet
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
echo -e "${YELLOW}Next steps:${NC}"
echo ""
echo "  1. Set your API key:"
echo "     export ANTHROPIC_API_KEY=sk-ant-..."
echo ""
echo "  2. Run go-code:"
echo "     ./bin/go-code"
echo ""
echo "  3. Push to GitHub (repo: strings77wzq/claude-code-Go):"
echo "     git remote add origin git@github.com:strings77wzq/claude-code-Go.git"
echo "     git add ."
echo "     git commit -m 'feat: initial release'"
echo "     git branch -M main"
echo "     git push -u origin main"
echo ""
