package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"mcp-bridge/internal/bridge"
)

// version is set by the build process
var version = "dev"

var (
	serverURL   = flag.String("server", "", "Remote MCP server URL (required)")
	apiKey      = flag.String("key", "", "API key for authentication (required)")
	debug       = flag.Bool("debug", false, "Enable all debug logging (equivalent to -debug-client -debug-server)")
	debugClient = flag.Bool("debug-client", false, "Enable client-side message logging")
	debugServer = flag.Bool("debug-server", false, "Enable server-side message logging")
	showVersion = flag.Bool("version", false, "Show version and exit")
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("mcp-bridge %s\n", version)
		os.Exit(0)
	}

	if *serverURL == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Create bridge with debug settings
	b := bridge.New(*serverURL, *apiKey, *debug)

	// Set granular debug flags (debug flag enables both)
	debugClientEnabled := *debug || *debugClient
	debugServerEnabled := *debug || *debugServer
	b.SetDebugFlags(debugClientEnabled, debugServerEnabled)

	if err := b.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
