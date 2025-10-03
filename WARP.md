# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

- Repository: mcp-bridge (Go)
- Purpose: Bridge stdio MCP communication to remote HTTP MCP servers, enabling transparent access to remote MCP capabilities.

Project summary
- Minimal Go CLI consisting of a single entrypoint (main.go) and a Go module (go.mod).
- Core behavior: acts as a local MCP server on stdio while proxying all communication to a remote HTTP-based MCP server using the official MCP Go SDK.

Prerequisites
- Go: >= 1.24.6 (from go.mod). Use your local Go toolchain (asdf/brew acceptable).
- Optional tooling: gitleaks (run before committing per house rules).

Quick start
- Build the binary, then start it with required flags to bridge stdio MCP to remote HTTP MCP.

Common commands
- Setup
  - Download module deps:
    ```bash path=null start=null
    go mod download
    ```
  - Keep go.mod/go.sum tidy:
    ```bash path=null start=null
    go mod tidy
    ```

- Build
  - Build a local binary:
    ```bash path=null start=null
    mkdir -p ./bin && go build -o ./bin/mcp-bridge .
    ```

- Run
  - Use environment variables for secrets, then pass flags. The program bridges stdio to HTTP MCP.
    ```bash path=null start=null
    export MCP_API_KEY={{MCP_API_KEY}}
    ./bin/mcp-bridge \
      -server "https://remote-mcp-server.com" \
      -key "$MCP_API_KEY" \
      -debug
    ```

- Tests
  - This repository uses BDD testing with godog. BDD tests are provided under bdd/ and features/.
    - Install BDD dependency and tidy:
      ```bash path=null start=null
      go get github.com/cucumber/godog@latest && go mod tidy
      ```
    - Run all BDD features (via go test):
      ```bash path=null start=null
      go test ./...
      ```
    - Run a single feature by name (regex matches scenario/test name):
      ```bash path=null start=null
      go test -run 'mcp-bridge-bdd' ./bdd -v
      ```

- Format and vet
  - Format source:
    ```bash path=null start=null
    go fmt ./...
    ```
  - Static checks:
    ```bash path=null start=null
    go vet ./...
    ```

- Secrets scan (house rule)
  - Scan repository for secrets before commit:
    ```bash path=null start=null
    gitleaks detect --source . --no-banner
    ```
  - Optional staged protection hook:
    ```bash path=null start=null
    gitleaks protect --staged --no-banner
    ```

High-level architecture
- Entry point: main.go defines flags (-server, -key, -debug) and creates a bridge instance.
- Core logic: internal/bridge package contains MCPBridge struct and MCP transport bridging logic.
- MCPBridge struct holds RemoteURL, APIKey, MCP server, MCP client, and Debug flag.
- Run loop (in internal/bridge):
  - Creates an MCP server that accepts stdio connections (left side of bridge).
  - Creates an MCP client that connects to remote HTTP server via SSE transport (right side).
  - Proxies all MCP protocol communication bidirectionally between stdio and HTTP.
  - Uses the official MCP Go SDK for protocol handling.
- Logging is gated behind -debug to avoid noisy output.

Files of interest
- go.mod — module name (mcp-bridge) and Go version.
- main.go — CLI entry point with flag parsing and bridge initialization.
- internal/bridge/ — core bridge logic separated from CLI concerns.

Extending the program
- Core logic lives in internal/bridge; add new functionality there.
- Add unit tests alongside internal packages (e.g., internal/bridge/bridge_test.go).
- BDD tests live in bdd/ and import internal/bridge.
- Keep main.go focused on flag parsing and wiring; prefer small, focused functions for easy unit testing.

House rules observed here
- Keep solutions simple and clean; underengineer until needed.
- Always run gitleaks before committing.
- If you introduce Rails or a Rails service in this repo in the future, use Rails generators to bootstrap new functionality.

Pre-commit hooks
- Pre-commit hooks are set up to enforce quality standards:
  - Install with: ./scripts/install-hooks.sh
  - Runs gitleaks, go vet, go fmt, BDD tests, and build verification
  - Prevents commits that fail any checks
- All contributors should install the hooks after cloning
