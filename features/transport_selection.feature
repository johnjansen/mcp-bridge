Feature: Bridge chooses the best available transport
  As a user running the bridge
  I want it to try streaming first and fall back to HTTP POST
  So that it works with both stream-capable and plain JSON-RPC servers

  Background:
    Given a remote MCP server base URL "http://localhost:3000"
    And debug logging is enabled

  Scenario: Streaming endpoint available
    Given the server exposes streaming at "/mcp/stream"
    And the server accepts streaming sessions
    When I start the bridge
    Then the bridge connects using streaming transport
    And the bridge logs "Using streaming transport"

  Scenario: Streaming endpoint not available, fallback succeeds
    Given the server does not expose streaming at "/mcp/stream"
    And the server accepts JSON-RPC over POST at "/mcp"
    When I start the bridge
    Then the bridge falls back to HTTP POST transport
    And the bridge logs "falling back to HTTP POST"
    And the bridge can process JSON-RPC requests

  Scenario: Both streaming and HTTP POST unavailable
    Given the server does not expose streaming at "/mcp/stream"
    And the server does not accept JSON-RPC over POST at "/mcp"
    When I start the bridge
    Then the bridge fails with an error about connection failed