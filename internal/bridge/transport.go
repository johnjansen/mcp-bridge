package bridge

import (
	"context"
	"io"
)

// Transport represents a bidirectional MCP connection to a remote endpoint
type Transport interface {
	io.ReadWriteCloser
	Connect(ctx context.Context) error
}

// Ensure compatibility with MCP SDK expectations
// We adapt our Transport to mcp.Transport inside bridge.go where needed.
