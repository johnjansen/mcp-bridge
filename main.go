package main

import (
	"flag"
	"log"
	"os"

	"mcp-bridge/internal/bridge"
)

var (
	serverURL = flag.String("server", "", "Remote MCP server URL (required)")
	apiKey    = flag.String("key", "", "API key for authentication (required)")
	debug     = flag.Bool("debug", false, "Enable debug logging")
)

func main() {
	flag.Parse()

	if *serverURL == "" || *apiKey == "" {
		flag.Usage()
		os.Exit(1)
	}

	b := bridge.New(*serverURL, *apiKey, *debug)
	if err := b.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
