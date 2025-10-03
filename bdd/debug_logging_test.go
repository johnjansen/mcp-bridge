package bdd

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/cucumber/godog"
	"mcp-bridge/internal/bridge"
)

type safeBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *safeBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

func (s *safeBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.String()
}

type debugWorld struct {
	remoteURL string
	apiKey    string
	bridge    *bridge.MCPBridge
	logBuffer *safeBuffer // Thread-safe buffer for concurrent log writes
}

func (w *debugWorld) reset() {
	w.logBuffer = &safeBuffer{}
	log.SetOutput(w.logBuffer) // Redirect log output to our buffer
}

func (w *debugWorld) aRemoteMCPServerAtWithAPIKey(url, key string) error {
	w.remoteURL = url
	w.apiKey = key
	return nil
}

func (w *debugWorld) anMCPBridgeWithoutAnyDebugFlags() error {
	w.bridge = bridge.New(w.remoteURL, w.apiKey, false)
	return nil
}

func (w *debugWorld) iStartTheBridgeWithFlag(flag string) error {
	var debug, debugClient, debugServer bool
	switch flag {
	case "-debug":
		debug = true
	case "-debug-client":
		debugClient = true
	case "-debug-server":
		debugServer = true
	}
	w.bridge = bridge.New(w.remoteURL, w.apiKey, debug)
	w.bridge.SetDebugFlags(debugClient, debugServer)

	// Start bridge in background
	go func() { _ = w.bridge.Run() }()
	// Add a small delay to let the bridge start
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (w *debugWorld) allDebugLoggingShouldBeEnabled() error {
	if !w.bridge.Debug || !w.bridge.DebugClient || !w.bridge.DebugServer {
		return fmt.Errorf("expected all debug flags to be enabled")
	}
	return nil
}

func (w *debugWorld) connectionLifecycleEventsShouldBeLogged() error {
	logs := w.logBuffer.String()
	if !strings.Contains(logs, "Starting MCP bridge") {
		return fmt.Errorf("expected connection lifecycle logs, got: %s", logs)
	}
	return nil
}

func (w *debugWorld) anMCPBridgeWithoutAnyDebugFlagsEnabled() error {
	w.bridge = bridge.New(w.remoteURL, w.apiKey, false)
	w.bridge.SetDebugFlags(false, false)
	return nil
}

func (w *debugWorld) aCompleteMCPExchangeOccurs() error {
	// Simulate complete MCP exchange
	// This will be implemented when we add the actual debug logging
	return nil
}

func (w *debugWorld) noDebugMessagesShouldBeLogged() error {
	logs := w.logBuffer.String()
	if strings.Contains(logs, "→") || strings.Contains(logs, "←") {
		return fmt.Errorf("unexpected debug logs found: %s", logs)
	}
	return nil
}

func (w *debugWorld) onlyErrorsShouldBeLogged() error {
	// In a real test, we'd verify that only error-level logs are present
	// For now, this is a placeholder
	return nil
}

func InitializeDebugScenario(ctx *godog.ScenarioContext) {
	world := &debugWorld{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		world.reset()
		return ctx, nil
	})

	ctx.Step(`^a remote MCP server at "([^"]*)" with API key "([^"]*)"$`, world.aRemoteMCPServerAtWithAPIKey)
	ctx.Step(`^an MCP bridge without any debug flags enabled$`, world.anMCPBridgeWithoutAnyDebugFlags)
	ctx.Step(`^I start the bridge with "([^"]*)" flag$`, world.iStartTheBridgeWithFlag)
	ctx.Step(`^all debug logging should be enabled$`, world.allDebugLoggingShouldBeEnabled)
	ctx.Step(`^connection lifecycle events should be logged$`, world.connectionLifecycleEventsShouldBeLogged)
	ctx.Step(`^an MCP bridge without any debug flags$`, world.anMCPBridgeWithoutAnyDebugFlagsEnabled)
	ctx.Step(`^a complete MCP exchange occurs$`, world.aCompleteMCPExchangeOccurs)
	ctx.Step(`^no debug messages should be logged$`, world.noDebugMessagesShouldBeLogged)
	ctx.Step(`^only errors should be logged$`, world.onlyErrorsShouldBeLogged)
}
