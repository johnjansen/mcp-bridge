#!/bin/bash

# Install pre-commit hooks for mcp-bridge

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
HOOKS_DIR="$REPO_ROOT/.git/hooks"

echo "🔧 Installing pre-commit hooks..."

# Check if we're in a git repository
if [ ! -d "$REPO_ROOT/.git" ]; then
    echo "❌ Error: Not in a git repository"
    exit 1
fi

# Check if gitleaks is installed
if ! command -v gitleaks &> /dev/null; then
    echo "⚠️  gitleaks not found. Install with:"
    echo "   brew install gitleaks"
    echo "   # or"
    echo "   go install github.com/gitleaks/gitleaks/v8@latest"
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Copy pre-commit hook
PRE_COMMIT_HOOK="$HOOKS_DIR/pre-commit"

if [ -f "$PRE_COMMIT_HOOK" ]; then
    echo "⚠️  Pre-commit hook already exists. Backing up..."
    mv "$PRE_COMMIT_HOOK" "$PRE_COMMIT_HOOK.backup.$(date +%s)"
fi

# Create the pre-commit hook
cat > "$PRE_COMMIT_HOOK" << 'EOF'
#!/bin/bash

set -e

echo "🔍 Running pre-commit checks..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_step() {
    echo -e "${YELLOW}🔧 $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if gitleaks is installed
if ! command -v gitleaks &> /dev/null; then
    print_error "gitleaks is not installed. Install with: brew install gitleaks"
    exit 1
fi

# 1. Run gitleaks secret scan
print_step "Scanning for secrets with gitleaks..."
if gitleaks detect --source . --no-banner; then
    print_success "No secrets detected"
else
    print_error "Secrets detected! Please remove them before committing."
    exit 1
fi

# 2. Run go vet
print_step "Running go vet..."
if go vet ./...; then
    print_success "go vet passed"
else
    print_error "go vet failed"
    exit 1
fi

# 3. Run go fmt check
print_step "Checking go fmt..."
UNFORMATTED=$(gofmt -l .)
if [ -z "$UNFORMATTED" ]; then
    print_success "All files are formatted"
else
    print_error "The following files are not formatted:"
    echo "$UNFORMATTED"
    echo "Run: go fmt ./..."
    exit 1
fi

# 4. Run BDD tests
print_step "Running BDD tests..."
if go test ./bdd -v; then
    print_success "BDD tests passed"
else
    print_error "BDD tests failed"
    exit 1
fi

# 5. Build check
print_step "Checking build..."
if go build -o /tmp/mcp-bridge .; then
    print_success "Build successful"
    rm -f /tmp/mcp-bridge
else
    print_error "Build failed"
    exit 1
fi

echo -e "${GREEN}🎉 All pre-commit checks passed!${NC}"
EOF

chmod +x "$PRE_COMMIT_HOOK"

echo "✅ Pre-commit hook installed successfully!"
echo ""
echo "The hook will run automatically before each commit and check:"
echo "  • Secrets scanning (gitleaks)"
echo "  • Code quality (go vet)"
echo "  • Code formatting (go fmt)"
echo "  • BDD tests"
echo "  • Build verification"
echo ""
echo "To test the hook manually, run:"
echo "  .git/hooks/pre-commit"