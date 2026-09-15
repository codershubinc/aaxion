#!/bin/bash
# Aaxion Build Script
# Compiles the project for multiple operating systems and architectures.

set -e

APP_NAME="aaxion"
BUILD_DIR="build"
MAIN_PATH="./cmd/aaxion"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔨 Building Aaxion for multiple platforms...${NC}"

# Clean previous build
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

# Ensure dependencies are tidy
go mod tidy

# List of OS and architectures to build for:
# Format: "OS/ARCH"
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "linux/arm"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for PLATFORM in "${PLATFORMS[@]}"; do
    OS=${PLATFORM%/*}
    ARCH=${PLATFORM#*/}
    
    OUTPUT_NAME="${APP_NAME}-${OS}-${ARCH}"
    
    # Add .exe extension for Windows
    if [ "$OS" = "windows" ]; then
        OUTPUT_NAME+=".exe"
    fi
    
    echo "   • Building for $OS/$ARCH -> $BUILD_DIR/$OUTPUT_NAME"
    
    # Run the build
    env GOOS=$OS GOARCH=$ARCH go build -ldflags="-s -w" -o "$BUILD_DIR/$OUTPUT_NAME" "$MAIN_PATH"
    
    if [ $? -ne 0 ]; then
        echo -e "${RED}❌ Failed to build for $OS/$ARCH${NC}"
        exit 1
    fi
done

echo -e "${GREEN}✅ All builds completed successfully! Assets are in the '$BUILD_DIR' directory.${NC}"
