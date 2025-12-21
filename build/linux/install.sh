#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Starting Aaxion Installation...${NC}"

# 1. Check for Root
if [ "$EUID" -ne 0 ]; then 
  echo -e "${RED}❌ Please run as root (sudo ./install.sh)${NC}"
  exit 1
fi

echo -e "${BLUE}Stopping service if running${NC}"
systemctl stop aaxion 2>/dev/null


# 2. Build the Project
echo -e "${BLUE}📦 Building Release Binary...${NC}"
# We assume the script is run from setup/linux, so we go up two levels to root
PROJECT_ROOT=$(dirname $(dirname $(dirname $(realpath $0))))
cd "$PROJECT_ROOT"

if ! cargo build --release; then
    echo -e "${RED}❌ Build failed! Please check your Rust installation.${NC}"
    exit 1
fi

# remove old installation if exists
echo -e "${BLUE}🧹 Removing old installation...${NC}
rm -rf /opt/aaxion"

# 3. Create Directories
echo -e "${BLUE}📂 Creating /opt/aaxion...${NC}"
mkdir -p /opt/aaxion
mkdir -p /opt/aaxion/uploads

# 4. Copy Files
echo -e "${BLUE}🚚 Copying files...${NC}"
cp target/release/aaxion /opt/aaxion/aaxion-server
cp -r assets /opt/aaxion/

# 5. Install Systemd Service
echo -e "${BLUE}⚙️  Installing Systemd Service...${NC}"
cp setup/linux-installer/aaxion.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable aaxion
systemctl restart aaxion

# 6. Install Desktop Shortcut (for the user who ran sudo)
# We get the real user who invoked sudo
REAL_USER=${SUDO_USER:-$USER}
USER_HOME=$(getent passwd "$REAL_USER" | cut -d: -f6)

echo -e "${BLUE}🖥️  Installing Desktop Shortcut for $REAL_USER...${NC}"
mkdir -p "$USER_HOME/.local/share/applications"
cp setup/linux-installer/aaxion.desktop "$USER_HOME/.local/share/applications/"
chown "$REAL_USER:$REAL_USER" "$USER_HOME/.local/share/applications/aaxion.desktop"

echo -e "${GREEN}✅ Installation Complete!${NC}"
echo -e "   - Service is running at: http://localhost:18875"
echo -e "   - Files are stored in:   /opt/aaxion/uploads"
echo -e "   - Check status with:     systemctl status aaxion"