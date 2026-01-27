#!/bin/sh
# Aaxion Installer Script for Linux and macOS

set -e

# --- Configuration ---
REPO="codershubinc/aaxion"
BINARY_NAME="aaxion"
# ---------------------

# 1. Detect OS & Architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Linux)  OS="linux" ;;
    Darwin) OS="darwin" ;;
    *)      echo "❌ Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *)      echo "❌ Unsupported Architecture: $ARCH"; exit 1 ;;
esac

echo "⚡️ Aaxion Installer"
echo "   • OS: $OS"
echo "   • Arch: $ARCH"

# 2. Find Latest Release (using GitHub API)
echo "🔍 Finding latest release..."
LATEST_URL=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep "browser_download_url" | grep "$OS" | grep "$ARCH" | cut -d '"' -f 4)

if [ -z "$LATEST_URL" ]; then
    echo "❌ Could not find a release for $OS/$ARCH."
    exit 1
fi

echo "   • Downloading from: $LATEST_URL"

# 3. Download & Install
TMP_DIR=$(mktemp -d)
curl -fsSL "$LATEST_URL" -o "$TMP_DIR/$BINARY_NAME"
chmod +x "$TMP_DIR/$BINARY_NAME"

# 4. Move to Path
INSTALL_DIR="/usr/local/bin"

if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
else
    echo "🔑 Sudo permission required to install to $INSTALL_DIR"
    sudo mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
fi

# 5. Cleanup & Verify
rm -rf "$TMP_DIR"

echo ""
echo "✅ Installed successfully!"
echo "   Run 'aaxion --help' to get started."