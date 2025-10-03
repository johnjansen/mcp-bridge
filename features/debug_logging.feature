Feature: Enhanced Debug Logging
  As a developer using MCP Bridge
  I want to monitor and debug MCP protocol interactions
  So that I can understand and troubleshoot communication between clients and servers

  Background:
    Given a remote MCP server at "https://example.com/mcp" with API key "test-key"
    And an MCP bridge without any debug flags enabled

  Scenario: Global debug logging enables all debug features
    When I start the bridge with "-debug" flag
    Then all debug logging should be enabled
    And connection lifecycle events should be logged

  Scenario: Client-side debug logging only
    When I start the bridge with "-debug-client" flag
    Then connection lifecycle events should be logged

  Scenario: Server-side debug logging only
    When I start the bridge with "-debug-server" flag
    Then connection lifecycle events should be logged

  Scenario: No debug logging when flags are disabled
    Given an MCP bridge without any debug flags
    When a complete MCP exchange occurs
    Then no debug messages should be logged
    And only errors should be logged
