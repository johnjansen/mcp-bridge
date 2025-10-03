#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print colored output
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}" >&2
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}" >&2
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}" >&2
}

print_error() {
    echo -e "${RED}❌ $1${NC}" >&2
}

# Configuration
REPO="johnjansen/mcp-bridge"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="mcp-bridge"

# Detect platform and architecture
detect_platform() {
    local os arch
    
    # Detect OS
    case "$(uname -s)" in
        Darwin*)
            os="darwin"
            ;;
        Linux*)
            os="linux"
            ;;
        *)
            print_error "Unsupported operating system: $(uname -s)"
            exit 1
            ;;
    esac
    
    # Detect architecture
    case "$(uname -m)" in
        x86_64|amd64)
            arch="amd64"
            ;;
        arm64|aarch64)
            arch="arm64"
            ;;
        *)
            print_error "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac
    
    echo "${os}_${arch}"
}

# Get the latest release version
get_latest_version() {
    local latest_url="https://api.github.com/repos/${REPO}/releases/latest"
    local version
    
    print_info "Fetching latest release information..."
    
    if command -v curl >/dev/null 2>&1; then
        version=$(curl -s "$latest_url" | grep '"tag_name":' | head -n 1 | cut -d '"' -f 4)
    elif command -v wget >/dev/null 2>&1; then
        version=$(wget -qO- "$latest_url" | grep '"tag_name":' | head -n 1 | cut -d '"' -f 4)
    else
        print_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi
    
    if [ -z "$version" ]; then
        print_error "Failed to fetch latest version"
        exit 1
    fi
    
    echo "$version"
}

# Download and install binary
install_binary() {
    local version="$1"
    local platform="$2"
    local tarball="${BINARY_NAME}-${version}-${platform}.tar.gz"
    local binary_name="${BINARY_NAME}-${platform}"
    local download_url="https://github.com/${REPO}/releases/download/${version}/${tarball}"
    local temp_dir=$(mktemp -d)
    
    print_info "Downloading ${BINARY_NAME} ${version} for ${platform}..."
    
    # Download tarball
    cd "$temp_dir"
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "$tarball" "$download_url"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "$tarball" "$download_url"
    else
        print_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi
    
    # Extract tarball
    print_info "Extracting binary..."
    tar -xzf "$tarball"
    
    # The extracted binary name includes the version
    local extracted_binary="${BINARY_NAME}-${version}-${platform}"
    
    # Install binary
    print_info "Installing to ${INSTALL_DIR}/${BINARY_NAME}..."
    
    # Check if we need sudo
    if [ -w "$INSTALL_DIR" ]; then
        mv "$extracted_binary" "${INSTALL_DIR}/${BINARY_NAME}"
    else
        print_warning "Requesting sudo access to install to ${INSTALL_DIR}..."
        sudo mv "$extracted_binary" "${INSTALL_DIR}/${BINARY_NAME}"
    fi
    
    # Make executable
    if [ -w "${INSTALL_DIR}/${BINARY_NAME}" ]; then
        chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    else
        sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    fi
    
    # Cleanup
    cd - >/dev/null
    rm -rf "$temp_dir"
    
    print_success "Successfully installed ${BINARY_NAME} ${version}"
}

# Verify installation
verify_installation() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local installed_version
        installed_version=$("$BINARY_NAME" -version 2>/dev/null | head -n 1)
        print_success "Installation verified: $installed_version"
        print_info "Run '${BINARY_NAME} -h' to see usage instructions"
    else
        print_error "Installation failed - ${BINARY_NAME} not found in PATH"
        exit 1
    fi
}

# Main installation process
main() {
    echo -e "${BLUE}"
    echo "🌉 MCP Bridge Installer"
    echo "======================="
    echo -e "${NC}"
    
    # Check for existing installation
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local current_version
        current_version=$("$BINARY_NAME" -version 2>/dev/null | awk '{print $2}' || echo "unknown")
        print_warning "${BINARY_NAME} is already installed (version: ${current_version})"
        read -p "Do you want to reinstall? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_info "Installation cancelled"
            exit 0
        fi
    fi
    
    # Detect platform
    local platform
    platform=$(detect_platform)
    print_info "Detected platform: $platform"
    
    # Get latest version
    local version
    version=$(get_latest_version)
    print_info "Latest version: $version"
    
    # Install binary
    install_binary "$version" "$platform"
    
    # Verify installation
    verify_installation
    
    echo
    print_success "🎉 Installation complete!"
    echo -e "${BLUE}Get started:${NC}"
    echo "  ${BINARY_NAME} -server \"https://your-mcp-server.com\" -key \"\$API_KEY\" -debug"
    echo
}

# Run main function
main "$@"