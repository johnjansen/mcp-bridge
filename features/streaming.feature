Feature: MCP Transport Bridge
  As a developer
  I want to use stdio to communicate with a remote MCP server
  So that my local MCP client can transparently access remote MCP capabilities

  Background:
    The bridge acts as a local MCP server on stdio while proxying to a remote HTTP MCP server
    This enables: stdio ↔ [stdio ↔ http] ↔ http

  Scenario: Bridge accepts stdio connections and connects to remote server
    Given a remote MCP server at "https://example.com/mcp" with API key "test-key"
    And an MCP bridge configured for that remote server
    When the bridge starts
    Then it should accept MCP connections on stdio
    And it should establish a connection to the remote server

  Scenario: Bridge proxies MCP initialize handshake
    Given a running MCP bridge connected to a remote server
    When a client sends an MCP initialize request via stdin
    Then the bridge forwards the initialize request to the remote server
    And the bridge returns the remote server's initialize response via stdout
    And the client-bridge-server connection is established

  Scenario: Bridge proxies tool listing
    Given an established MCP bridge connection
    When a client requests the list of available tools
    Then the bridge forwards the tools/list request to the remote server
    And returns the remote server's tool list to the client

  Scenario: Bridge proxies tool execution
    Given an established MCP bridge connection
    When a client calls a tool with arguments
    Then the bridge forwards the tools/call request to the remote server
    And returns the tool execution result to the client

  Scenario: Bridge handles bidirectional communication
    Given an established MCP bridge connection
    When the remote server sends a notification to the bridge
    Then the bridge forwards the notification to the client via stdout
    And when the client sends a request, it's forwarded to the remote server
