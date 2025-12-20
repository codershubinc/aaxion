#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 1. Navigate to Project Root
# We assume the script is run from build/linux, so we go up two levels to root
PROJECT_ROOT=$(dirname $(dirname $(dirname $(realpath "$0"))))
cd "$PROJECT_ROOT" || { echo -e "${RED}❌ Failed to navigate to project root${NC}"; exit 1; }

echo -e "${BLUE}🏗️  Starting Build Process...${NC}"

# 2. Build Release Binary
if cargo build --release; then
    echo -e "${GREEN}✅ Build Successful${NC}"
else
    echo -e "${RED}❌ Build Failed! Check Rust installation.${NC}"
    exit 1
fi

# 3. Prepare Installer Directory
INSTALLER_DIR="setup/linux-installer"
echo -e "${BLUE}🧹 Cleaning up old installer files...${NC}"
rm -rf "$INSTALLER_DIR"
mkdir -p "$INSTALLER_DIR"

# 4. Copy Files
echo -e "${BLUE}📦 Packaging files...${NC}"

# Binary (Rename to aaxion-server)
cp target/release/localdrive-rs "$INSTALLER_DIR/aaxion-server" || { echo "❌ Binary copy failed"; exit 1; }

# Assets (CRITICAL: The web UI needs these)
cp -r assets "$INSTALLER_DIR/" || { echo "❌ Assets copy failed"; exit 1; }

# Config/Service Files
# Note: Adjust these paths if your source files are in setup/linux/ or build/linux/
# Based on your previous context, they seem to be in setup/linux/
cp build/linux/aaxion.desktop "$INSTALLER_DIR/"
cp build/linux/aaxion.service "$INSTALLER_DIR/"
cp build/linux/install.sh "$INSTALLER_DIR/"

# Make installer executable
chmod +x "$INSTALLER_DIR/install.sh"

echo -e "${GREEN}✅ Package Created at: $INSTALLER_DIR${NC}"
ls -l "$INSTALLER_DIR"

# 5. Create Tarball (Optional but recommended for distribution)
echo -e "${BLUE}🎁 Creating compressed archive...${NC}"
tar -czf aaxion-linux-installer.tar.gz -C setup linux-installer
echo -e "${GREEN}🎉 Ready to ship: aaxion-linux-installer.tar.gz${NC}"

echo -e "\nTo test locally:"
echo -e "  cd $INSTALLER_DIR"
echo -e "  sudo ./install.sh"