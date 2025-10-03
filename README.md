# MCP Bridge

[![CI](https://github.com/johnjansen/mcp-bridge/workflows/CI/badge.svg)](https://github.com/johnjansen/mcp-bridge/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/johnjansen/mcp-bridge)](https://goreportcard.com/report/github.com/johnjansen/mcp-bridge)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/johnjansen/mcp-bridge)](https://golang.org/)

> **🚀 Stream Any Data to Your MCP Server in Real-Time**
>
> Turn any command-line tool into an MCP data source. Pipe logs, metrics, sensor data, or API responses directly to your Model Context Protocol server with built-in authentication and channel organization. Perfect for feeding AI agents with live data streams, building monitoring pipelines, or creating custom data integrations.
>
> **Why MCP Bridge?**
> - **Zero Configuration**: Just pipe data and go
> - **Secure by Design**: Bearer token authentication built-in
> - **Channel Organization**: Route different data types to specific channels
> - **Bulletproof Streaming**: Handles backpressure and connection issues gracefully
> - **Universal Input**: Works with any tool that writes to stdout

## Installation

Download the latest release or build from source:

```bash
# Build from source
git clone https://github.com/johnjansen/mcp-bridge.git
cd mcp-bridge
go build -o bin/mcp-bridge .
```

## Quick Start

Stream data to your MCP server in three simple steps:

```bash
# 1. Set your API key (never hardcode secrets)
export MCP_API_KEY="your-secret-api-key"

# 2. Pipe any data to the bridge
echo "Hello, MCP!" | ./bin/mcp-bridge \
  -server "https://your-mcp-server.com" \
  -key "$MCP_API_KEY" \
  -channel "alerts"

# 3. Your data is now streaming to your MCP server
```

## Real-World Examples

### Stream Application Logs
```bash
tail -f /var/log/app.log | ./bin/mcp-bridge \
  -server "https://mcp.example.com" \
  -key "$MCP_API_KEY" \
  -channel "app-logs"
```

### Monitor System Metrics
```bash
vmstat 1 | ./bin/mcp-bridge \
  -server "https://mcp.example.com" \
  -key "$MCP_API_KEY" \
  -channel "system-metrics"
```

### Stream API Data
```bash
curl -N "https://api.example.com/events" | ./bin/mcp-bridge \
  -server "https://mcp.example.com" \
  -key "$MCP_API_KEY" \
  -channel "api-events"
```

### Database Change Stream
```bash
pg_receivexlog --slot=changes --stdout | ./bin/mcp-bridge \
  -server "https://mcp.example.com" \
  -key "$MCP_API_KEY" \
  -channel "db-changes"
```

## Command Line Options

| Flag | Description | Required |
|------|-------------|----------|
| `-server` | MCP server URL (e.g., `https://mcp.example.com`) | Yes |
| `-key` | API key for authentication | Yes |
| `-channel` | Channel name for organizing data streams | Yes |
| `-debug` | Enable debug logging | No |

## How It Works

MCP Bridge reads data from stdin in 4KB chunks and posts each chunk to your MCP server via HTTP:

```
stdin → [4KB buffer] → POST /api/v1/stream/{channel}
```

**HTTP Details:**
- **Method**: POST
- **URL**: `{server}/api/v1/stream/{channel}`
- **Headers**: 
  - `Authorization: Bearer {key}`
  - `Content-Type: application/octet-stream`
- **Body**: Raw binary data from stdin

## Security

- **Never hardcode API keys** in commands or scripts
- Always use environment variables: `export MCP_API_KEY="your-key"`
- Use HTTPS URLs for your MCP server
- Rotate API keys regularly

## Development

This project uses BDD testing with Gherkin scenarios:

```bash
# Install dependencies
go get github.com/cucumber/godog@latest
go mod tidy

# Run all tests
go test ./...

# Run BDD tests with verbose output
go test -run 'mcp-bridge-bdd' ./bdd -v

# Build the binary
mkdir -p ./bin && go build -o ./bin/mcp-bridge .

# Format and vet
go fmt ./...
go vet ./...
```

### Project Structure

```
├── main.go                 # CLI entry point
├── internal/bridge/        # Core bridge logic
│   ├── bridge.go          # MCPBridge struct and HTTP client
│   └── runner.go          # Streaming loop
├── bdd/                   # BDD tests
│   ├── steps_test.go      # Godog step definitions
│   └── suite_test.go      # Test suite runner
├── features/              # Gherkin scenarios
│   └── streaming.feature  # Core streaming behavior
└── WARP.md               # Development guidelines
```

## Architecture

- **Minimal Design**: Single-purpose tool that does one thing well
- **Streaming First**: Handles continuous data streams efficiently
- **Error Resilient**: Goroutines with proper error propagation
- **Memory Efficient**: 4KB buffer with bounded memory usage
- **Debug Support**: Optional verbose logging for troubleshooting

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

This tool bridges the gap between any command-line data source and MCP servers. It's the missing piece that lets you stream real-time data to AI agents and other MCP consumers without complex integrations.