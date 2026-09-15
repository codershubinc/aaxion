#!/bin/sh
# Aaxion Installer Script

set -e

# --- Configuration ---
REPO="codershubinc/aaxion"
BINARY_NAME="aaxion"
# ---------------------

# Colors
CYAN='\033[0;36m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "${BLUE}=======================================${NC}"
echo "${CYAN}      Aaxion Installation Script       ${NC}"
echo "${BLUE}=======================================${NC}"
echo ""

# 1. Fetch Latest Version
echo "${CYAN}[1/4]${NC} Fetching latest version info..."
LATEST_VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    # Fallback if the user hasn't published a formal release yet (e.g. beta release not marked as latest)
    LATEST_VERSION="latest"
    echo "      ${YELLOW}Note: Could not determine tag from GitHub API. Defaulting to 'latest' release branch.${NC}"
else
    echo "      ${GREEN}Found version: ${LATEST_VERSION}${NC}"
fi

# 2. Detect OS & Architecture
echo "${CYAN}[2/4]${NC} Detecting system architecture..."
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    armv7l) ARCH="arm" ;;
    *) echo "${RED}Error: Unsupported Architecture: $ARCH${NC}"; exit 1 ;;
esac

echo "      OS:   ${OS}"
echo "      Arch: ${ARCH}"

# 3. Download
ASSET_NAME="${BINARY_NAME}-${OS}-${ARCH}"
if [ "$LATEST_VERSION" = "latest" ]; then
    DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${ASSET_NAME}"
else
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${ASSET_NAME}"
fi

TMP_DIR=$(mktemp -d)
echo "${CYAN}[3/4]${NC} Downloading ${ASSET_NAME}..."

# Use curl with progress bar (-#) and follow redirects (-L) while failing on errors (-f)
curl -# -fL "$DOWNLOAD_URL" -o "$TMP_DIR/$BINARY_NAME" || {
    echo "${RED}Error: Failed to download release asset from ${DOWNLOAD_URL}${NC}"
    rm -rf "$TMP_DIR"
    exit 1
}

chmod +x "$TMP_DIR/$BINARY_NAME"

# 4. Install
INSTALL_DIR="/usr/local/bin"
FINAL_NAME="aax"

echo "${CYAN}[4/4]${NC} Installing to $INSTALL_DIR/$FINAL_NAME..."

if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$FINAL_NAME"
else
    echo "      ${YELLOW}Sudo privileges required to write to $INSTALL_DIR${NC}"
    sudo mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$FINAL_NAME"
fi

# Cleanup
rm -rf "$TMP_DIR"

if command -v $FINAL_NAME >/dev/null 2>&1; then
    echo ""
    echo "${GREEN}Installation successful!${NC}"
    echo "Run '${CYAN}aax --help${NC}' to get started."
else
    echo ""
    echo "${RED}Warning: Installation finished, but $FINAL_NAME is not in your PATH.${NC}"
fi
