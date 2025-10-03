Feature: Stream stdin to MCP server
  As a developer
  I want the bridge to POST stdin chunks to the MCP server with auth
  So that downstream consumers receive the data by channel

  Scenario: Streams a single chunk to the configured channel
    Given a test MCP server
    And a bridge configured for that server with api key "test-key" and channel "test-channel"
    When I stream the message "hello"
    Then the server received 1 request to path "/api/v1/stream/test-channel"
    And the request had header "Authorization" equal to "Bearer test-key"
    And the last request body equals "hello"
