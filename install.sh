#!/bin/sh
# Aaxion Installer Script

set -e

# --- Configuration ---
REPO="codershubinc/aaxion"
BINARY_NAME="aaxion"
FINAL_NAME="aax"
INSTALL_DIR="/usr/local/bin"
# ---------------------

# --- Setup Colors ---
CYAN=$(printf '\033[0;36m')
GREEN=$(printf '\033[0;32m')
BLUE=$(printf '\033[0;34m')
RED=$(printf '\033[0;31m')
YELLOW=$(printf '\033[1;33m')
BOLD=$(printf '\033[1m')
NC=$(printf '\033[0m')

# --- Helper Functions ---
step() { echo "${CYAN}${BOLD}[$1/4]${NC} $2"; }
success() { echo "      ${GREEN}✔${NC} $1"; }
info() { echo "      ${BLUE}ℹ${NC} $1"; }
warn() { echo "      ${YELLOW}⚠${NC} $1"; }
error() { echo "      ${RED}✖ Error:${NC} $1"; exit 1; }

# --- Ensure Cleanup on Exit ---
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT INT TERM

# --- Header ---
echo ""
echo "${BLUE}${BOLD}=======================================${NC}"
echo "${CYAN}${BOLD}      🚀 Aaxion Installation Script    ${NC}"
echo "${BLUE}${BOLD}=======================================${NC}"
echo ""

# --- Pre-flight Checks ---
command -v curl >/dev/null 2>&1 || error "'curl' is required but not installed."
command -v grep >/dev/null 2>&1 || error "'grep' is required but not installed."

# 1. Fetch Latest Version
step "1" "Fetching latest version info..."
LATEST_VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    LATEST_VERSION="latest"
    warn "Could not determine tag from GitHub. Defaulting to 'latest'."
else
    success "Found version: ${BOLD}${LATEST_VERSION}${NC}"
fi

# 2. Detect OS & Architecture
step "2" "Detecting system architecture..."
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    armv7l) ARCH="arm" ;;
    *) error "Unsupported Architecture: $ARCH" ;;
esac

info "OS:   ${OS}"
info "Arch: ${ARCH}"

# 3. Download
ASSET_NAME="${BINARY_NAME}-${OS}-${ARCH}"
if [ "$LATEST_VERSION" = "latest" ]; then
    DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${ASSET_NAME}"
else
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${ASSET_NAME}"
fi

step "3" "Downloading ${ASSET_NAME}..."
curl -# -fL "$DOWNLOAD_URL" -o "$TMP_DIR/$BINARY_NAME" || error "Failed to download release asset."
chmod +x "$TMP_DIR/$BINARY_NAME"
success "Download complete."

# 4. Install
step "4" "Installing to ${INSTALL_DIR}/${FINAL_NAME}..."

if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$FINAL_NAME"
else
    warn "Sudo privileges required to write to ${INSTALL_DIR}"
    sudo mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$FINAL_NAME"
fi

success "Successfully installed ${FINAL_NAME}."
echo ""

# --- Final Check ---
if command -v $FINAL_NAME >/dev/null 2>&1; then
    echo "${GREEN}${BOLD}🎉 Installation successful!${NC}"
    echo "Run '${CYAN}${BOLD}${FINAL_NAME} --help${NC}' to get started."
else
    echo "${YELLOW}${BOLD}⚠ Installation finished, but ${FINAL_NAME} is not in your PATH.${NC}"
    echo "Please add ${INSTALL_DIR} to your environment PATH."
fi
echo ""