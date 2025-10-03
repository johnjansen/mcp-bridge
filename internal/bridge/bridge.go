package bridge

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

type MCPBridge struct {
	ServerURL string
	APIKey    string
	Channel   string
	Client    *http.Client
	Debug     bool
}

func New(serverURL, apiKey, channel string, debug bool) *MCPBridge {
	return &MCPBridge{
		ServerURL: serverURL,
		APIKey:    apiKey,
		Channel:   channel,
		Client:    &http.Client{},
		Debug:     debug,
	}
}

func (b *MCPBridge) Log(format string, v ...interface{}) {
	if b.Debug {
		log.Printf(format, v...)
	}
}

func (b *MCPBridge) StreamToServer(data []byte) error {
	url := fmt.Sprintf("%s/api/v1/stream/%s", b.ServerURL, b.Channel)
	
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", b.APIKey))
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := b.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned error: %s", resp.Status)
	}

	b.Log("Successfully sent %d bytes to server", len(data))
	return nil
}