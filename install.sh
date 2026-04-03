#!/bin/bash
# install.sh — One-command installer for go-code
# Usage: curl -fsSL https://raw.githubusercontent.com/strings77wzq/claude-code-Go/main/install.sh | bash
#
# Detects OS/arch, downloads pre-built binary, installs to /usr/local/bin.

set -e

REPO="strings77wzq/claude-code-Go"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="go-code"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  go-code Installer${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Detect OS
case "$(uname -s)" in
    Linux*)  OS="linux" ;;
    Darwin*) OS="darwin" ;;
    *)       echo -e "${RED}Unsupported OS: $(uname -s)${NC}"; exit 1 ;;
esac

# Detect arch
case "$(uname -m)" in
    x86_64)  ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    arm64)   ARCH="arm64" ;;
    *)       echo -e "${RED}Unsupported arch: $(uname -m)${NC}"; exit 1 ;;
esac

BINARY="${BINARY_NAME}-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY}"

echo -e "${YELLOW}Detected: ${OS}/${ARCH}${NC}"
echo -e "${YELLOW}Downloading: ${BINARY}${NC}"

# Download
if command -v curl &> /dev/null; then
    curl -fsSL -o "/tmp/${BINARY}" "${DOWNLOAD_URL}"
elif command -v wget &> /dev/null; then
    wget -q -O "/tmp/${BINARY}" "${DOWNLOAD_URL}"
else
    echo -e "${RED}Error: curl or wget required${NC}"
    exit 1
fi

# Install
chmod +x "/tmp/${BINARY}"
if [ -w "${INSTALL_DIR}" ]; then
    mv "/tmp/${BINARY}" "${INSTALL_DIR}/${BINARY_NAME}"
else
    echo -e "${YELLOW}Need sudo to install to ${INSTALL_DIR}${NC}"
    sudo mv "/tmp/${BINARY}" "${INSTALL_DIR}/${BINARY_NAME}"
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Installation Complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "Binary installed to: ${INSTALL_DIR}/${BINARY_NAME}"
echo ""
echo "Next steps:"
echo "  1. Set your API key:"
echo "     export ANTHROPIC_API_KEY=sk-ant-..."
echo ""
echo "  2. Start using go-code:"
echo "     go-code"
echo ""
