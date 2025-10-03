package bridge

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// httpPostTransport implements a simple HTTP POST request/response transport
// This mimics the Ruby bridge behavior for compatibility with servers that don't support streaming
type httpPostTransport struct {
	endpoint   string
	httpClient *http.Client
	debug      bool
}

func newHTTPPostTransport(endpoint string, client *http.Client, debug bool) *httpPostTransport {
	if client == nil {
		client = http.DefaultClient
	}
	return &httpPostTransport{
		endpoint:   endpoint,
		httpClient: client,
		debug:      debug,
	}
}

// Run starts the HTTP POST bridge loop
func (t *httpPostTransport) Run(ctx context.Context) error {
	if t.debug {
		log.Printf("HTTP POST bridge running, reading from stdin...")
	}

	stdin := bufio.NewReader(os.Stdin)
	stdout := os.Stdout

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Read from stdin
			line, err := stdin.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					if t.debug {
						log.Printf("EOF on stdin, shutting down")
					}
					return nil
				}
				return fmt.Errorf("read error: %w", err)
			}

			data := bytes.TrimSpace(line)
			if len(data) == 0 {
				continue
			}

			// Validate it's JSON
			var msg map[string]interface{}
			if err := json.Unmarshal(data, &msg); err != nil {
				log.Printf("Invalid JSON received: %v", err)
				continue
			}

			// Send to remote server via HTTP POST
			if t.debug {
				log.Printf("→ Sending to %s: %s", t.endpoint, string(data))
			}

			req, err := http.NewRequestWithContext(ctx, "POST", t.endpoint, bytes.NewReader(data))
			if err != nil {
				log.Printf("Failed to create request: %v", err)
				continue
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			resp, err := t.httpClient.Do(req)
			if err != nil {
				log.Printf("Request failed: %v", err)
				// Create error response
				errorResp := map[string]interface{}{
					"jsonrpc": "2.0",
					"id":      msg["id"],
					"error": map[string]interface{}{
						"code":    -32603,
						"message": fmt.Sprintf("Bridge error: %v", err),
					},
				}
				errorData, _ := json.Marshal(errorResp)
				stdout.Write(errorData)
				stdout.Write([]byte("\n"))
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Printf("HTTP %d: %s", resp.StatusCode, resp.Status)
				// Create error response
				errorResp := map[string]interface{}{
					"jsonrpc": "2.0",
					"id":      msg["id"],
					"error": map[string]interface{}{
						"code":    -32603,
						"message": fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status),
					},
				}
				errorData, _ := json.Marshal(errorResp)
				stdout.Write(errorData)
				stdout.Write([]byte("\n"))
				continue
			}

			// Read response
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Failed to read response: %v", err)
				continue
			}

			if t.debug {
				log.Printf("← Received: %s", string(body))
			}

			// Write response to stdout
			stdout.Write(body)
			stdout.Write([]byte("\n"))
		}
	}
}
