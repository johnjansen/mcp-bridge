# MCP Bridge

[![CI](https://github.com/johnjansen/mcp-bridge/workflows/CI/badge.svg)](https://github.com/johnjansen/mcp-bridge/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/johnjansen/mcp-bridge)](https://goreportcard.com/report/github.com/johnjansen/mcp-bridge)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/johnjansen/mcp-bridge)](https://golang.org/)

> **🌉 Bridge stdio to Remote MCP Servers with Built-in Observability**
>
> A drop-in bridge with bidirectional logging that lets you:
> - **Watch MCP Traffic** → See requests and responses with directional indicators (→/←)
> - **Debug Selectively** → Monitor just client or server side with granular controls
> - **Troubleshoot Fast** → Capture minimal traces without changing client code
> - **Learn & Demo** → Perfect for workshops, demos, and teaching MCP flows
>
> Plus traditional bridging: `stdio ↔ [stdio ↔ http] ↔ http`
>
> Even if your tools support HTTP streaming directly, mcp-bridge adds observability without changing your workflow.

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
While Cursor and VSCode support HTTP streaming for MCP directly, you can still use mcp-bridge to intercept and debug traffic:
```json
{
  "mcpServers": {
    "your-mcp-server-debug": {
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
}
```
This configuration lets you observe all MCP traffic in real-time, perfect for debugging and understanding protocol interactions.

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
| `-debug` | Enable all debug logging | No |
| `-debug-client` | Enable client-side message logging | No |
| `-debug-server` | Enable server-side message logging | No |

### Debug Logging

MCP Bridge provides granular debug logging to help troubleshoot communication issues and observe MCP traffic in real time. This makes it a powerful companion even when your IDE can already use HTTP streaming for MCP directly.

```bash
# Full debug logging (both client and server)
mcp-bridge -server "https://example.com/mcp" -key "$API_KEY" -debug

# Client-side only (→ indicator for client messages)
mcp-bridge -server "https://example.com/mcp" -key "$API_KEY" -debug-client

# Server-side only (← indicator for server messages)
mcp-bridge -server "https://example.com/mcp" -key "$API_KEY" -debug-server
```

Example debug output:
```
2025/10/03 17:40:43 Starting MCP bridge to https://example.com/mcp (debug: global=true, client=true, server=true)
→ Initialize request from client:
{
  "version": "1.0.0",
  "protocol": "mcp"
}
→ Forwarding to remote server
← Response from remote server:
{
  "ok": true,
  "version": "1.0.0"
}
← Forwarding to client
```

#### Why use mcp-bridge even when you “don’t really need” a bridge?

Even if your IDE or tool can talk to your remote MCP server directly, running traffic through mcp-bridge gives you:

- Observable, directional logs
  - See each request and response with clear direction markers: → client-to-server, ← server-to-client
  - Turn on only the side you care about using -debug-client or -debug-server
- Faster troubleshooting and support
  - Capture minimal, anonymized traces to reproduce issues without exposing full payloads
  - Spot protocol mismatches and malformed messages early
- Non-invasive monitoring in dev/staging
  - Drop-in in front of any existing MCP client without changing client code
  - Keep your existing workflow while gaining visibility
- Teaching and demos
  - Show newcomers what “MCP over stdio/http” actually looks like
  - Great for workshops, screen shares, and bug bashes

At-a-glance examples:

```bash
# Watch just the client side (requests going out)
mcp-bridge -server "https://example.com/mcp" -key "$API_KEY" -debug-client

# Watch just the server side (responses and notifications coming back)
mcp-bridge -server "https://example.com/mcp" -key "$API_KEY" -debug-server

# Full-duplex tracing for a short session to capture an issue
mcp-bridge -server "https://example.com/mcp" -key "$API_KEY" -debug | tee trace.log
```

## How It Works

MCP Bridge creates a bidirectional proxy between stdio and HTTP:

```
Local MCP Client ←→ stdio ←→ [MCP Bridge] ←→ HTTP ←→ Remote MCP Server
```

### Transport Mechanisms

MCP Bridge supports two transport mechanisms for communicating with remote servers, automatically selecting the best available option:

1. **MCP Streaming Transport** (Primary)
   - Native MCP streaming protocol over HTTP
   - Connects to `/stream` endpoint for efficient message exchange
   - Uses the official MCP Go SDK's StreamableClientTransport
   - Ideal for real-time tool execution and notifications
   - Maintains continuous connection with the server

2. **HTTP POST Transport** (Fallback)
   - Traditional request-response communication
   - Each MCP message sent as separate HTTP POST request
   - Compatible with servers that don't support streaming
   - Implements JSON-RPC over HTTP protocol
   - Mimics Ruby bridge behavior for maximum compatibility

**Transport Selection Process:**
1. Bridge attempts streaming connection to `/stream` endpoint
2. Uses 3-second timeout to test streaming capability
3. If streaming succeeds, establishes bidirectional HTTP streaming transport
4. If streaming fails or times out, falls back to HTTP POST
5. Logs transport selection when debug enabled

Example debug output during transport negotiation:
```
2025/10/03 17:40:43 Attempting streaming transport...
2025/10/03 17:40:43 Using streaming transport
```
or with fallback:
```
2025/10/03 17:40:43 Attempting streaming transport...
2025/10/03 17:40:46 Streaming not supported (connection refused), falling back to HTTP POST
```

**Protocol Details:**
- **Input**: JSON-RPC MCP protocol via stdin/stdout
- **Output**: HTTP requests to remote MCP server (streaming or POST)
- **Authentication**: Bearer token authentication
- **Bidirectional**: Handles both client requests and server notifications
- **Content Type**: application/json with chunked transfer encoding

### Transport Configuration Examples

**1. Server with Both Transports:**
No special configuration needed - bridge auto-negotiates:
```bash
mcp-bridge -server "https://mcp.example.com" -key "$API_KEY" -debug
```

**2. Streaming-Only Server:**
Ensure `/stream` endpoint is available:
```bash
mcp-bridge -server "https://streaming.example.com" -key "$API_KEY" -debug
# Bridge will use HTTP streaming transport exclusively
```

**3. Legacy Server (HTTP POST only):**
Bridge automatically falls back to HTTP POST:
```bash
mcp-bridge -server "https://legacy.example.com" -key "$API_KEY" -debug
# After 3s timeout, falls back to HTTP POST transport
```

**Debug Output Examples:**

1. Successful Streaming Connection:
```
2025/10/03 17:40:43 Starting MCP bridge to https://example.com/mcp (debug: global=true)
2025/10/03 17:40:43 Attempting streaming transport...
2025/10/03 17:40:43 Using streaming transport
2025/10/03 17:40:43 Connected to remote MCP server
```

2. Fallback to HTTP POST:
```
2025/10/03 17:40:43 Starting MCP bridge to https://legacy.example.com/mcp (debug: global=true)
2025/10/03 17:40:43 Attempting streaming transport...
2025/10/03 17:40:46 Streaming not supported (connection refused), falling back to HTTP POST
2025/10/03 17:40:46 HTTP POST bridge running, reading from stdin...
```

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