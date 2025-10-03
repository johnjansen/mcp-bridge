Feature: HTTP POST transport request/response behavior
  As a user running the bridge in fallback mode
  I want it to handle JSON-RPC requests and responses correctly
  So that MCP communication works reliably

  Background:
    Given a remote MCP server base URL "http://localhost:3000"
    And debug logging is enabled
    And the server does not expose streaming at "/mcp/stream"
    And the server accepts JSON-RPC over POST at "/mcp"
    And I start the bridge

  Scenario: Successful initialize request
    When I send a JSON-RPC request:
      """
      {
        "jsonrpc": "2.0",
        "method": "initialize",
        "id": 1,
        "params": {
          "protocolVersion": "0.1.0",
          "capabilities": {},
          "clientInfo": {
            "name": "test",
            "version": "1.0"
          }
        }
      }
      """
    Then I receive a JSON-RPC response with id 1 and result {}

  Scenario: Parse error from server
    When I send malformed JSON:
      """
      {not valid json
      """
    Then I receive a JSON-RPC error response with code -32700
    And the error message contains "Parse error"

  Scenario: Method not found error
    When I send a JSON-RPC request:
      """
      {
        "jsonrpc": "2.0",
        "method": "nonexistent_method",
        "id": 2,
        "params": {}
      }
      """
    Then I receive a JSON-RPC error response with code -32601
    And the error message contains "Method not found"

  Scenario: Bridge handles concurrent requests
    When I send 5 concurrent JSON-RPC "ping" requests
    Then I receive 5 successful responses with matching IDs