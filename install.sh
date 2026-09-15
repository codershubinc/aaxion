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

printf "%b\n" "${BLUE}=======================================${NC}"
printf "%b\n" "${CYAN}      Aaxion Installation Script       ${NC}"
printf "%b\n" "${BLUE}=======================================${NC}"
echo ""

# 1. Fetch Latest Version
printf "%b\n" "${CYAN}[1/4]${NC} Fetching latest version info..."
LATEST_VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    # Fallback if the user hasn't published a formal release yet (e.g. beta release not marked as latest)
    LATEST_VERSION="latest"
    printf "%b\n" "      ${YELLOW}Note: Could not determine tag from GitHub API. Defaulting to 'latest' release branch.${NC}"
else
    printf "%b\n" "      ${GREEN}Found version: ${LATEST_VERSION}${NC}"
fi

# 2. Detect OS & Architecture
printf "%b\n" "${CYAN}[2/4]${NC} Detecting system architecture..."
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    armv7l) ARCH="arm" ;;
    *) printf "%b\n" "${RED}Error: Unsupported Architecture: $ARCH${NC}"; exit 1 ;;
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
printf "%b\n" "${CYAN}[3/4]${NC} Downloading ${ASSET_NAME}..."

# Use curl with progress bar (-#) and follow redirects (-L) while failing on errors (-f)
curl -# -fL "$DOWNLOAD_URL" -o "$TMP_DIR/$BINARY_NAME" || {
    printf "%b\n" "${RED}Error: Failed to download release asset from ${DOWNLOAD_URL}${NC}"
    rm -rf "$TMP_DIR"
    exit 1
}

chmod +x "$TMP_DIR/$BINARY_NAME"

# 4. Install
INSTALL_DIR="/usr/local/bin"
FINAL_NAME="aax"

printf "%b\n" "${CYAN}[4/4]${NC} Installing to $INSTALL_DIR/$FINAL_NAME..."

if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$FINAL_NAME"
else
    printf "%b\n" "      ${YELLOW}Sudo privileges required to write to $INSTALL_DIR${NC}"
    sudo mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$FINAL_NAME"
fi

# Cleanup
rm -rf "$TMP_DIR"

if command -v $FINAL_NAME >/dev/null 2>&1; then
    echo ""
    printf "%b\n" "${GREEN}Installation successful!${NC}"
    printf "%b\n" "Run '${CYAN}aax --help${NC}' to get started."
else
    echo ""
    printf "%b\n" "${RED}Warning: Installation finished, but $FINAL_NAME is not in your PATH.${NC}"
fi
