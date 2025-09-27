#!/bin/bash

# Build script for GophKeeper
# This script builds the client and server for multiple platforms

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Build configuration
VERSION=${VERSION:-"1.0.0"}
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILD_DIR="build"
LDFLAGS="-X main.version=${VERSION} -X main.buildDate=${BUILD_DATE}"

# Platforms to build for
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
    "windows/arm64"
    "darwin/amd64"
    "darwin/arm64"
)

echo -e "${GREEN}Building GophKeeper v${VERSION}${NC}"
echo -e "${YELLOW}Build date: ${BUILD_DATE}${NC}"

# Create build directory
mkdir -p ${BUILD_DIR}

# Build server
echo -e "${GREEN}Building server...${NC}"
go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/gophkeeper-server ./cmd/server

# Build client for each platform
echo -e "${GREEN}Building client for multiple platforms...${NC}"
for platform in "${PLATFORMS[@]}"; do
    IFS='/' read -r os arch <<< "${platform}"
    
    echo -e "${YELLOW}Building for ${os}/${arch}...${NC}"
    
    # Set environment variables for cross-compilation
    export GOOS=${os}
    export GOARCH=${arch}
    
    # Build client
    output_name="gophkeeper-client-${os}-${arch}"
    if [ "${os}" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${output_name} ./cmd/client
    
    echo -e "${GREEN}✓ Built ${output_name}${NC}"
done

# Reset environment variables
unset GOOS
unset GOARCH

# Create archive for each platform
echo -e "${GREEN}Creating archives...${NC}"
cd ${BUILD_DIR}

for platform in "${PLATFORMS[@]}"; do
    IFS='/' read -r os arch <<< "${platform}"
    
    client_name="gophkeeper-client-${os}-${arch}"
    if [ "${os}" = "windows" ]; then
        client_name="${client_name}.exe"
    fi
    
    archive_name="gophkeeper-${os}-${arch}-v${VERSION}.tar.gz"
    
    if [ "${os}" = "windows" ]; then
        # Create ZIP for Windows
        archive_name="gophkeeper-${os}-${arch}-v${VERSION}.zip"
        zip -q ${archive_name} ${client_name} gophkeeper-server
    else
        # Create TAR.GZ for Unix-like systems
        tar -czf ${archive_name} ${client_name} gophkeeper-server
    fi
    
    echo -e "${GREEN}✓ Created ${archive_name}${NC}"
done

cd ..

echo -e "${GREEN}Build completed successfully!${NC}"
echo -e "${YELLOW}Build artifacts are in the ${BUILD_DIR}/ directory${NC}"
