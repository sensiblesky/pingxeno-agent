#!/bin/bash

set -e

VERSION=${1:-"dev"}
BUILD_DIR="build"
MAIN_PACKAGE="./cmd/agent"

# Build flags
LDFLAGS="-X main.version=${VERSION}"

# Create build directory
mkdir -p ${BUILD_DIR}

# Build for multiple platforms
echo "Building PingXeno Agent..."

# Linux AMD64
echo "Building for linux/amd64..."
GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/pingxeno-agent-linux-amd64 ${MAIN_PACKAGE}

# Linux ARM64
echo "Building for linux/arm64..."
GOOS=linux GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/pingxeno-agent-linux-arm64 ${MAIN_PACKAGE}

# Windows AMD64
echo "Building for windows/amd64..."
GOOS=windows GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/pingxeno-agent-windows-amd64.exe ${MAIN_PACKAGE}

# macOS AMD64
echo "Building for darwin/amd64..."
GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/pingxeno-agent-darwin-amd64 ${MAIN_PACKAGE}

# macOS ARM64 (Apple Silicon)
echo "Building for darwin/arm64..."
GOOS=darwin GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/pingxeno-agent-darwin-arm64 ${MAIN_PACKAGE}

# FreeBSD AMD64
echo "Building for freebsd/amd64..."
GOOS=freebsd GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/pingxeno-agent-freebsd-amd64 ${MAIN_PACKAGE}

echo "Build complete! Binaries are in ${BUILD_DIR}/"

