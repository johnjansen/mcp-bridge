package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	serverURL = flag.String("server", "", "MCP server URL (required)")
	apiKey    = flag.String("key", "", "API key for authentication (required)")
	channel   = flag.String("channel", "", "Channel name (required)")
	debug     = flag.Bool("debug", false, "Enable debug logging")
)

type MCPBridge struct {
	serverURL string
	apiKey    string
	channel   string
	client    *http.Client
	debug     bool
}

func NewMCPBridge(serverURL, apiKey, channel string, debug bool) *MCPBridge {
	return &MCPBridge{
		serverURL: serverURL,
		apiKey:    apiKey,
		channel:   channel,
		client:    &http.Client{},
		debug:     debug,
	}
}

func (b *MCPBridge) log(format string, v ...interface{}) {
	if b.debug {
		log.Printf(format, v...)
	}
}

func (b *MCPBridge) streamToServer(data []byte) error {
	url := fmt.Sprintf("%s/api/v1/stream/%s", b.serverURL, b.channel)
	
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", b.apiKey))
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := b.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned error: %s", resp.Status)
	}

	b.log("Successfully sent %d bytes to server", len(data))
	return nil
}

func (b *MCPBridge) Run() error {
	reader := bufio.NewReader(os.Stdin)
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		buffer := make([]byte, 4096)

		for {
			n, err := reader.Read(buffer)
			if err == io.EOF {
				b.log("Reached EOF, stopping")
				return
			}
			if err != nil {
				errChan <- fmt.Errorf("error reading from stdin: %v", err)
				return
			}

			if err := b.streamToServer(buffer[:n]); err != nil {
				errChan <- fmt.Errorf("error streaming to server: %v", err)
				return
			}
		}
	}()

	// Wait for either completion or error
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Return first error if any
	if err := <-errChan; err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()

	if *serverURL == "" || *apiKey == "" || *channel == "" {
		flag.Usage()
		os.Exit(1)
	}

	bridge := NewMCPBridge(*serverURL, *apiKey, *channel, *debug)
	
	if err := bridge.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
