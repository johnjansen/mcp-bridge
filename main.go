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
	debug       = flag.Bool("debug", false, "Enable debug logging")
	showVersion = flag.Bool("version", false, "Show version and exit")
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("mcp-bridge %s\n", version)
		os.Exit(0)
	}

	if *serverURL == "" || *apiKey == "" {
		flag.Usage()
		os.Exit(1)
	}

	b := bridge.New(*serverURL, *apiKey, *debug)
	if err := b.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
