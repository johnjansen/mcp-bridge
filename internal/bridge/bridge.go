package bridge

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCPBridge manages bidirectional communication between a local stdio MCP client
// and a remote HTTP MCP server.
type MCPBridge struct {
	RemoteURL   string
	APIKey      string
	Debug       bool // Global debug flag (enables all debugging)
	DebugClient bool // Enable client-side message logging
	DebugServer bool // Enable server-side message logging
	server      *mcp.Server
	client      *mcp.Client
	ctx         context.Context
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

// SetDebugFlags configures granular debug logging flags
func (b *MCPBridge) SetDebugFlags(debugClient, debugServer bool) {
	b.DebugClient = b.Debug || debugClient
	b.DebugServer = b.Debug || debugServer
}

// LogClient logs client-side messages with → indicator
func (b *MCPBridge) LogClient(format string, v ...interface{}) {
	if b.Debug || b.DebugClient {
		log.Printf("→ "+format, v...)
	}
}

// LogServer logs server-side messages with ← indicator
func (b *MCPBridge) LogServer(format string, v ...interface{}) {
	if b.Debug || b.DebugServer {
		log.Printf("← "+format, v...)
	}
}

// Log logs general messages (not specific to client/server)
func (b *MCPBridge) Log(format string, v ...interface{}) {
	if b.Debug || b.DebugClient || b.DebugServer {
		log.Printf(format, v...)
	}
}

// formatMCPMessage formats an MCP protocol message for logging
func (b *MCPBridge) formatMCPMessage(msg interface{}) string {
	if msg == nil {
		return "<nil>"
	}

	// Try to marshal the message to JSON for readability
	json, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		return fmt.Sprintf("%+v", msg)
	}
	return string(json)
}

// LogMCPClient logs client-side MCP protocol messages
func (b *MCPBridge) LogMCPClient(desc string, msg interface{}) {
	if b.Debug || b.DebugClient {
		log.Printf("→ %s:\n%s", desc, b.formatMCPMessage(msg))
	}
}

// LogMCPServer logs server-side MCP protocol messages
func (b *MCPBridge) LogMCPServer(desc string, msg interface{}) {
	if b.Debug || b.DebugServer {
		log.Printf("← %s:\n%s", desc, b.formatMCPMessage(msg))
	}
}

func (b *MCPBridge) Run() error {
	b.Log("Starting MCP bridge to %s (debug: global=%v, client=%v, server=%v)",
		b.RemoteURL, b.Debug, b.DebugClient, b.DebugServer)

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
