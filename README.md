# MCP Bridge

[![CI](https://github.com/johnjansen/mcp-bridge/workflows/CI/badge.svg)](https://github.com/johnjansen/mcp-bridge/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/johnjansen/mcp-bridge)](https://goreportcard.com/report/github.com/johnjansen/mcp-bridge)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/johnjansen/mcp-bridge)](https://golang.org/)

> **🌉 Bridge stdio to Remote MCP Servers**
>
> Connect your local MCP clients to remote MCP servers over HTTP. MCP Bridge acts as a local MCP server on stdio while proxying all communication to a remote HTTP-based MCP server, enabling transparent access to remote MCP capabilities.
>
> **Architecture**: `stdio ↔ [stdio ↔ http] ↔ http`

## Installation

### Automated Installer (Recommended)

The easiest way to install `mcp-bridge` is using our automated installer that detects your platform and downloads the appropriate binary:

```bash
curl -fsSL https://raw.githubusercontent.com/johnjansen/mcp-bridge/main/install.sh | bash
```

Or if you prefer to inspect the script first:
```bash
curl -fsSL https://raw.githubusercontent.com/johnjansen/mcp-bridge/main/install.sh -o install.sh
chmod +x install.sh
./install.sh
```

### Manual Installation from Releases

Download the latest release for your platform from the [releases page](https://github.com/johnjansen/mcp-bridge/releases):

```bash
# Linux AMD64
curl -L -o mcp-bridge.tar.gz https://github.com/johnjansen/mcp-bridge/releases/latest/download/mcp-bridge-latest-linux_amd64.tar.gz

# Linux ARM64
curl -L -o mcp-bridge.tar.gz https://github.com/johnjansen/mcp-bridge/releases/latest/download/mcp-bridge-latest-linux_arm64.tar.gz

# macOS Intel
curl -L -o mcp-bridge.tar.gz https://github.com/johnjansen/mcp-bridge/releases/latest/download/mcp-bridge-latest-darwin_amd64.tar.gz

# macOS Apple Silicon
curl -L -o mcp-bridge.tar.gz https://github.com/johnjansen/mcp-bridge/releases/latest/download/mcp-bridge-latest-darwin_arm64.tar.gz

# Extract and install
tar -xzf mcp-bridge.tar.gz
sudo mv mcp-bridge-* /usr/local/bin/mcp-bridge
chmod +x /usr/local/bin/mcp-bridge
```

### Build from Source

```bash
git clone https://github.com/johnjansen/mcp-bridge.git
cd mcp-bridge
go build -o bin/mcp-bridge .
```

## Quick Start

Connect to a remote MCP server in three simple steps:

```bash
# 1. Set your API key (never hardcode secrets)
export MCP_API_KEY="your-secret-api-key"

# 2. Start the bridge
mcp-bridge \
  -server "https://your-remote-mcp-server.com" \
  -key "$MCP_API_KEY" \
  -debug

# 3. Your local MCP client can now communicate via stdio
```

## IDE Configuration

MCP Bridge can be configured in various IDEs to enable local AI capabilities:

### Warp.dev
```json
{
  "your-mcp-server": {
    "command": "mcp-bridge",
    "args": [
      "-server", "https://your-remote-mcp-server.com",
      "-key", "$MCP_API_KEY",
      "-debug"
    ],
    "env": {
      "MCP_API_KEY": "your-secret-api-key"
    },
    "working_directory": null
  }
}
```

### Zed.dev
```json
{
  "your-mcp-server": {
    "command": "mcp-bridge",
    "args": [
      "-server", "https://your-remote-mcp-server.com",
      "-key", "$MCP_API_KEY",
      "-debug"
    ],
    "env": {
      "MCP_API_KEY": "your-secret-api-key"
    }
  }
}
```

### Cursor/VSCode
Cursor and VSCode support HTTP/SSE MCP servers directly, so you typically don't need mcp-bridge. Configure them directly:
```json
{
  "mcpServers": {
    "your-mcp-server": {
      "url": "https://your-remote-mcp-server.com/mcp",
      "headers": {
        "Authorization": "Bearer your-secret-api-key"
      }
    }
  }
}
```

## Use Cases

### Connect Local AI Tools to Remote MCP Servers
```bash
# Bridge enables local tools to access remote MCP capabilities
my-local-mcp-client | mcp-bridge -server "https://remote-mcp.com" -key "$API_KEY"
```

### Access Remote Tools and Resources
The bridge transparently proxies:
- **MCP initialize handshake** - Establishes connection
- **Tool listing and execution** - Access remote tools
- **Resource access** - Read remote resources  
- **Bidirectional communication** - Server notifications and client requests

### Development and Testing
```bash
# Connect to local development MCP server
mcp-bridge -server "http://localhost:3000" -key "dev-key" -debug
```

## Command Line Options

| Flag | Description | Required |
|------|-------------|----------|
| `-server` | Remote MCP server URL (HTTP/HTTPS) | Yes |
| `-key` | API key for authentication | Yes |
| `-debug` | Enable debug logging | No |

## How It Works

MCP Bridge creates a bidirectional proxy between stdio and HTTP:

```
Local MCP Client ←→ stdio ←→ [MCP Bridge] ←→ HTTP ←→ Remote MCP Server
```

**Protocol Details:**
- **Input**: JSON-RPC MCP protocol via stdin/stdout
- **Output**: HTTP requests to remote MCP server (SSE transport)
- **Authentication**: Bearer token authentication
- **Bidirectional**: Handles both client requests and server notifications

## Development

This project uses BDD testing with Gherkin scenarios:

```bash
# Install dependencies
go mod tidy

# Run all tests
go test ./...

# Run BDD tests with verbose output
go test ./bdd -v

# Build the binary
go build -o bin/mcp-bridge .

# Format and vet
go fmt ./...
go vet ./...
```

### Pre-commit Hooks

This project includes pre-commit hooks that automatically run:
- **gitleaks** - Secret scanning
- **go vet** - Static analysis
- **go fmt** - Code formatting check
- **BDD tests** - All scenarios must pass
- **Build verification** - Ensures code compiles

```bash
# Install pre-commit hooks
./scripts/install-hooks.sh

# Test hooks manually
.git/hooks/pre-commit
```

**Note**: Install gitleaks first: `brew install gitleaks`

### Project Structure

```
├── main.go                 # CLI entry point
├── internal/bridge/        # Core bridge logic
│   └── bridge.go          # MCP transport bridge
├── bdd/                   # BDD tests
│   ├── steps_test.go      # Godog step definitions
│   └── suite_test.go      # Test suite runner
├── features/              # Gherkin scenarios
│   └── streaming.feature  # MCP transport bridge behavior
└── WARP.md               # Development guidelines
```

## Architecture

- **Transport Bridge**: Bridges stdio MCP to HTTP MCP protocols
- **Bidirectional**: Handles both client-to-server and server-to-client communication
- **Authentication**: Secure API key authentication with remote servers
- **Error Resilient**: Proper error handling and connection management
- **Standards Compliant**: Uses official MCP Go SDK

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add BDD tests for new functionality
4. Ensure all tests pass: `go test ./...`
5. Run gitleaks before committing: `gitleaks detect --source .`
6. Submit a pull request

## License

MIT License - see LICENSE file for details.

## Why "Bridge"?

This tool bridges the gap between local stdio-based MCP clients and remote HTTP-based MCP servers. It enables transparent communication across different MCP transport protocols, making remote MCP capabilities accessible to local tools.