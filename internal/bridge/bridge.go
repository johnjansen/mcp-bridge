package bridge

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type MCPBridge struct {
	RemoteURL string
	APIKey    string
	Debug     bool
	server    *mcp.Server
	client    *mcp.Client
	ctx       context.Context
}

func New(remoteURL, apiKey string, debug bool) *MCPBridge {
	ctx := context.Background()

	// Create a server that accepts stdio connections (left side)
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "mcp-bridge",
		Version: "v1.0.0",
	}, nil)

	// Create a client to connect to remote server (right side)
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "mcp-bridge-client",
		Version: "v1.0.0",
	}, nil)

	return &MCPBridge{
		RemoteURL: remoteURL,
		APIKey:    apiKey,
		Debug:     debug,
		server:    server,
		client:    client,
		ctx:       ctx,
	}
}

func (b *MCPBridge) Log(format string, v ...interface{}) {
	if b.Debug {
		log.Printf(format, v...)
	}
}

func (b *MCPBridge) Run() error {
	b.Log("Starting MCP bridge to %s", b.RemoteURL)

	// Parse remote URL to determine transport type
	remoteURL, err := url.Parse(b.RemoteURL)
	if err != nil {
		return fmt.Errorf("invalid remote URL: %v", err)
	}

	// Connect to remote MCP server
	var transport mcp.Transport
	switch remoteURL.Scheme {
	case "http", "https":
		// Use SSE transport for remote HTTP connection
		client := &http.Client{}
		// Note: Authorization headers will need to be handled at the HTTP client level
		transport = &mcp.SSEClientTransport{
			Endpoint:   b.RemoteURL,
			HTTPClient: client,
		}
	default:
		return fmt.Errorf("unsupported URL scheme: %s", remoteURL.Scheme)
	}

	// Connect client to remote server
	remoteSession, err := b.client.Connect(b.ctx, transport, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to remote MCP server: %v", err)
	}
	defer remoteSession.Close()

	b.Log("Connected to remote MCP server")

	// Set up proxy server that forwards all requests to remote
	b.setupProxyHandlers(remoteSession)

	// Run the stdio server (this blocks)
	stdioTransport := &mcp.StdioTransport{}
	return b.server.Run(b.ctx, stdioTransport)
}

func (b *MCPBridge) setupProxyHandlers(remoteSession *mcp.ClientSession) {
	// This is where we'd set up proxy handlers for different MCP methods
	// For now, this is a placeholder - the MCP SDK might need custom handlers
	// or we might need to implement the proxy at a lower level
	b.Log("Setting up proxy handlers for remote session")

	// TODO: Implement actual proxy logic
	// This might require:
	// 1. Intercepting all incoming MCP requests from stdio
	// 2. Forwarding them to remoteSession
	// 3. Returning responses back to stdio client
	// 4. Handling bidirectional communication (server-initiated messages)
}
