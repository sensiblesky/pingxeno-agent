#!/bin/bash

# PingXeno Agent - Multi-OS Build Script
# This script builds the agent for all supported operating systems and architectures

set -e

VERSION=${1:-"1.0.0"}
BUILD_DIR="builds"
BINARY_NAME="pingxeno-agent"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Building PingXeno Agent v${VERSION} for all platforms...${NC}\n"

# Create build directory
mkdir -p ${BUILD_DIR}

# Build targets: OS/ARCH pairs
declare -a TARGETS=(
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
    "windows/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "freebsd/amd64"
    "freebsd/arm64"
)

# Function to get binary extension
get_ext() {
    local os=$1
    if [ "$os" = "windows" ]; then
        echo ".exe"
    else
        echo ""
    fi
}

# Function to get archive format
get_archive() {
    local os=$1
    if [ "$os" = "windows" ]; then
        echo "zip"
    else
        echo "tar.gz"
    fi
}

# Build for each target
for target in "${TARGETS[@]}"; do
    IFS='/' read -r os arch <<< "$target"
    ext=$(get_ext "$os")
    archive=$(get_archive "$os")
    
    echo -e "${YELLOW}Building for ${os}/${arch}...${NC}"
    
    # Set environment variables for cross-compilation
    export GOOS=$os
    export GOARCH=$arch
    
    # Build binary
    output_name="${BINARY_NAME}${ext}"
    output_path="${BUILD_DIR}/${BINARY_NAME}-${os}-${arch}${ext}"
    
    go build -ldflags="-s -w -X main.version=${VERSION}" -o "${output_path}" ./cmd/agent
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Built: ${output_path}${NC}"
        
        # Create archive
        archive_name="${BUILD_DIR}/${BINARY_NAME}-${os}-${arch}-${VERSION}.${archive}"
        
        if [ "$archive" = "zip" ]; then
            zip -j "${archive_name}" "${output_path}" > /dev/null
        else
            tar -czf "${archive_name}" -C "${BUILD_DIR}" "${BINARY_NAME}-${os}-${arch}${ext}" 2>/dev/null
        fi
        
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✓ Archived: ${archive_name}${NC}"
        fi
    else
        echo -e "${RED}✗ Failed to build for ${os}/${arch}${NC}"
    fi
    
    echo ""
done

# Reset GOOS and GOARCH
unset GOOS
unset GOARCH

echo -e "${GREEN}Build complete! Binaries are in the ${BUILD_DIR}/ directory${NC}"

