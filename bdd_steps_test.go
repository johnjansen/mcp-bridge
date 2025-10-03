package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/cucumber/godog"
)

type capturedRequest struct {
	path   string
	headers http.Header
	body   []byte
}

type world struct {
	server   *httptest.Server
	bridge   *MCPBridge
	mu       sync.Mutex
	requests []capturedRequest
}

func (w *world) iHaveATestMCPServer() error {
	w.requests = nil
	h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		_ = r.Body.Close()
		w.mu.Lock()
		w.requests = append(w.requests, capturedRequest{path: r.URL.Path, headers: r.Header.Clone(), body: append([]byte(nil), b...)})
		w.mu.Unlock()
		rw.WriteHeader(http.StatusOK)
	})
	w.server = httptest.NewServer(h)
	return nil
}

func (w *world) aBridgeConfiguredForThatServerWithKeyAndChannel(key, channel string) error {
	if w.server == nil {
		return fmt.Errorf("server not initialized")
	}
	w.bridge = NewMCPBridge(w.server.URL, key, channel, true)
	return nil
}

func (w *world) iStreamTheMessage(msg string) error {
	if w.bridge == nil {
		return fmt.Errorf("bridge not initialized")
	}
	return w.bridge.streamToServer([]byte(msg))
}

func (w *world) theServerReceivedNRequestsToPath(n int, path string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.requests) != n {
		return fmt.Errorf("expected %d requests, got %d", n, len(w.requests))
	}
	for i, req := range w.requests {
		if req.path != path {
			return fmt.Errorf("request %d path mismatch: expected %q got %q", i, path, req.path)
		}
	}
	return nil
}

func (w *world) theRequestHadHeaderEqualTo(name, expected string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.requests) == 0 {
		return fmt.Errorf("no requests captured")
	}
	last := w.requests[len(w.requests)-1]
	vals := last.headers.Values(name)
	if len(vals) == 0 {
		return fmt.Errorf("header %q missing", name)
	}
	for _, v := range vals {
		if v == expected {
			return nil
		}
	}
	return fmt.Errorf("header %q did not match %q, got %v", name, expected, vals)
}

func (w *world) theLastRequestBodyEquals(expected string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.requests) == 0 {
		return fmt.Errorf("no requests captured")
	}
	last := w.requests[len(w.requests)-1]
	if !bytes.Equal(last.body, []byte(expected)) {
		return fmt.Errorf("body mismatch: expected %q got %q", expected, string(last.body))
	}
	return nil
}

func (w *world) cleanup() {
	if w.server != nil {
		w.server.Close()
	}
}

func InitializeScenario(sc *godog.ScenarioContext) {
	w := &world{}
	sc.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) { return ctx, nil })
	sc.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) { w.cleanup(); return ctx, nil })

	sc.Step(`^a test MCP server$`, w.iHaveATestMCPServer)
	sc.Step(`^a bridge configured for that server with api key "([^"]*)" and channel "([^"]*)"$`, w.aBridgeConfiguredForThatServerWithKeyAndChannel)
	sc.Step(`^I stream the message "([^"]*)"$`, w.iStreamTheMessage)
	sc.Step(`^the server received (\d+) request to path "([^"]*)"$`, w.theServerReceivedNRequestsToPath)
	sc.Step(`^the request had header "([^"]*)" equal to "([^"]*)"$`, w.theRequestHadHeaderEqualTo)
	sc.Step(`^the last request body equals "([^"]*)"$`, w.theLastRequestBodyEquals)
}
