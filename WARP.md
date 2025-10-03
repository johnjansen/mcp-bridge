# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

- Repository: mcp-bridge (Go)
- Purpose: Stream stdin data to a remote MCP server over HTTP with bearer auth, organized by channel.

Project summary
- Minimal Go CLI consisting of a single entrypoint (main.go) and a Go module (go.mod).
- Core behavior: reads from stdin in chunks and POSTs each chunk to {SERVER}/api/v1/stream/{CHANNEL} with Authorization: Bearer {API_KEY}.

Prerequisites
- Go: >= 1.24.6 (from go.mod). Use your local Go toolchain (asdf/brew acceptable).
- Optional tooling: gitleaks (run before committing per house rules).

Quick start
- Build the binary, then pipe data to it with required flags.

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
  - Use environment variables for secrets, then pass flags. The program reads from stdin.
    ```bash path=null start=null
    export MCP_API_KEY={{MCP_API_KEY}}
    echo "hello world" | ./bin/mcp-bridge \
      -server "https://example.com" \
      -key "$MCP_API_KEY" \
      -channel "my-channel" \
      -debug
    ```

- Tests
  - This repository currently contains no unit tests. A BDD scaffold using godog is provided under bdd/ and features/.
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
      go test -run 'mcp-bridge-bdd' -v
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
- Entry point: main.go defines flags (-server, -key, -channel, -debug), constructs MCPBridge, and runs it.
- MCPBridge struct holds serverURL, apiKey, channel, an http.Client, and debug flag.
- Run loop:
  - Reads from stdin via bufio.Reader into a 4KB buffer in a goroutine.
  - For each non-empty read, posts bytes to {serverURL}/api/v1/stream/{channel}.
  - Uses Authorization: Bearer {apiKey}, Content-Type: application/octet-stream.
  - sync.WaitGroup and an error channel coordinate completion and propagate the first error.
- Logging is gated behind -debug to avoid noisy output.

Files of interest
- go.mod — module name (mcp-bridge) and Go version.
- main.go — entire implementation (stdin streaming and HTTP POST logic).

Extending the program
- If functionality grows, consider moving logic into internal/ packages (e.g., internal/bridge) and keeping main.go focused on flag parsing and wiring.
- Add _test.go files alongside new packages; prefer small, focused functions for easy unit testing.

House rules observed here
- Keep solutions simple and clean; underengineer until needed.
- Always run gitleaks before committing.
- If you introduce Rails or a Rails service in this repo in the future, use Rails generators to bootstrap new functionality.
