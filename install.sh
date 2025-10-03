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

# Detect best install directory from PATH
detect_install_dir() {
    local os="$(uname -s)"
    local arch="$(uname -m)"
    local candidate_dirs=()
    
    # OS-specific directory preferences
    case "$os" in
        Darwin*)
            # macOS: prefer Homebrew locations, then system, then user
            if [ "$arch" = "arm64" ]; then
                # Apple Silicon
                candidate_dirs=(
                    "/opt/homebrew/bin"
                    "/usr/local/bin"
                    "$HOME/.local/bin"
                    "$HOME/bin"
                )
            else
                # Intel Mac
                candidate_dirs=(
                    "/usr/local/bin"
                    "/opt/homebrew/bin"
                    "$HOME/.local/bin"
                    "$HOME/bin"
                )
            fi
            ;;
        Linux*)
            # Linux: prefer user-local (no sudo needed), then system
            candidate_dirs=(
                "$HOME/.local/bin"
                "$HOME/bin"
                "/usr/local/bin"
            )
            ;;
    esac
    
    # First, check if any candidate dir exists and is in PATH
    for dir in "${candidate_dirs[@]}"; do
        if [[ ":$PATH:" == *":$dir:"* ]]; then
            if [ -d "$dir" ] && [ -w "$dir" ]; then
                echo "$dir"
                return 0
            elif [ -d "$dir" ]; then
                # Directory exists but not writable - may need sudo
                echo "$dir"
                return 0
            fi
        fi
    done
    
    # If no existing dir found, create appropriate fallback
    local fallback
    case "$os" in
        Darwin*)
            # On macOS, try to create /usr/local/bin if possible (more standard)
            fallback="/usr/local/bin"
            if [ ! -d "$fallback" ]; then
                print_info "Creating $fallback..."
                sudo mkdir -p "$fallback"
                sudo chown $(whoami):admin "$fallback" 2>/dev/null || true
            fi
            ;;
        Linux*)
            # On Linux, use ~/.local/bin (FreeDesktop standard)
            fallback="$HOME/.local/bin"
            if [ ! -d "$fallback" ]; then
                print_info "Creating $fallback..."
                mkdir -p "$fallback"
            fi
            ;;
    esac
    
    # Check if it's in PATH
    if [[ ":$PATH:" != *":$fallback:"* ]]; then
        print_warning "$fallback is not in your PATH"
        case "$os" in
            Darwin*)
                print_info "Add this line to your ~/.zshrc:"
                ;;
            Linux*)
                print_info "Add this line to your ~/.bashrc or ~/.zshrc:"
                ;;
        esac
        echo "    export PATH=\"$fallback:\$PATH\""
    fi
    
    echo "$fallback"
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
    local install_dir="$3"
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
    
    # The extracted binary name is just binary-platform (no version)
    local extracted_binary="${BINARY_NAME}-${platform}"
    
    # Install binary
    print_info "Installing to ${install_dir}/${BINARY_NAME}..."
    
    # Create install directory if it doesn't exist
    if [ ! -d "$install_dir" ]; then
        print_info "Creating ${install_dir}..."
        if [[ "$install_dir" == "$HOME"* ]]; then
            mkdir -p "$install_dir"
        else
            sudo mkdir -p "$install_dir"
        fi
    fi
    
    # Check if we need sudo
    if [ -w "$install_dir" ]; then
        mv "$extracted_binary" "${install_dir}/${BINARY_NAME}"
    else
        print_warning "Requesting sudo access to install to ${install_dir}..."
        sudo mv "$extracted_binary" "${install_dir}/${BINARY_NAME}"
    fi
    
    # Make executable
    if [ -w "${install_dir}/${BINARY_NAME}" ]; then
        chmod +x "${install_dir}/${BINARY_NAME}"
    else
        sudo chmod +x "${install_dir}/${BINARY_NAME}"
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
    
    # Detect install directory
    local install_dir
    install_dir=$(detect_install_dir)
    print_info "Install directory: $install_dir"
    
    # Get latest version
    local version
    version=$(get_latest_version)
    print_info "Latest version: $version"
    
    # Install binary
    install_binary "$version" "$platform" "$install_dir"
    
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