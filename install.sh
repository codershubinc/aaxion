#!/bin/sh
# Aaxion Installer Script for Linux and macOS

set -e

# --- Configuration ---
REPO="codershubinc/aaxion"
BINARY_NAME="aaxion"
# ---------------------

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}⚡️ Aaxion Installer${NC}"

# 1. Detect OS & Architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    armv7l) ARCH="arm" ;;
    *) echo -e "${RED}❌ Unsupported Architecture: $ARCH${NC}"; exit 1 ;;
esac

echo "   • OS: $OS"
echo "   • Arch: $ARCH"

# 2. Build the Direct Download URL
ASSET_NAME="${BINARY_NAME}-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${ASSET_NAME}"

echo "🔍 Downloading latest release..."
echo "   • URL: $DOWNLOAD_URL"

# 3. Download & Install
TMP_DIR=$(mktemp -d)

# Use curl to download, capturing the HTTP status code
HTTP_STATUS=$(curl -fsSL -w "%{http_code}" "$DOWNLOAD_URL" -o "$TMP_DIR/$BINARY_NAME")

if [ "$HTTP_STATUS" != "200" ] && [ "$HTTP_STATUS" != "302" ]; then
    echo -e "${RED}❌ Failed to download release. Are you sure a release asset named '${ASSET_NAME}' exists?${NC}"
    rm -rf "$TMP_DIR"
    exit 1
fi

chmod +x "$TMP_DIR/$BINARY_NAME"

# 4. Move to Path
INSTALL_DIR="/usr/local/bin"
FINAL_NAME="aax"

echo "📦 Installing to $INSTALL_DIR as '$FINAL_NAME'..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$FINAL_NAME"
else
    echo "🔑 Sudo permission required to install to $INSTALL_DIR"
    sudo mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$FINAL_NAME"
fi

# 5. Cleanup & Verify
rm -rf "$TMP_DIR"

if command -v $FINAL_NAME >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Installed successfully!${NC}"
    echo "   Run '$FINAL_NAME --help' to see commands."
else
    echo -e "${RED}❌ Installation finished, but $FINAL_NAME is not in your PATH.${NC}"
fi
