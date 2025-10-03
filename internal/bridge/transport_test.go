package bridge

import (
	"context"
	"io"
	"os/exec"
	"strings"
	"testing"
)

// mockProcess implements a simple echo process for testing
type mockProcess struct {
	in  io.WriteCloser
	out io.ReadCloser
	err io.ReadCloser
	cmd *exec.Cmd
}

func newMockProcess() (*mockProcess, error) {
	// Use the 'cat' command as a mock echo server
	cmd := exec.Command("cat")

	in, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	out, err := cmd.StdoutPipe()
	if err != nil {
		in.Close()
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		in.Close()
		out.Close()
		return nil, err
	}

	return &mockProcess{in, out, stderr, cmd}, nil
}

func (p *mockProcess) start() error {
	return p.cmd.Start()
}

func (p *mockProcess) close() error {
	p.in.Close()
	p.out.Close()
	p.err.Close()
	return p.cmd.Wait()
}

func TestStdioTransport(t *testing.T) {
	ctx := context.Background()

	t.Run("echo process communication", func(t *testing.T) {
		// Create and start a mock process
		proc, err := newMockProcess()
		if err != nil {
			t.Fatalf("Failed to create mock process: %v", err)
		}
		if err := proc.start(); err != nil {
			t.Fatalf("Failed to start mock process: %v", err)
		}
		defer proc.close()

		// Create transport
		transport := &StdioTransport{
			cmd:    proc.cmd,
			stdin:  proc.in,
			stdout: proc.out,
			stderr: proc.err,
			debug:  true,
		}

		// Test communication
		message := "hello world\n"
		n, err := transport.Write([]byte(message))
		if err != nil {
			t.Fatalf("Write failed: %v", err)
		}
		if n != len(message) {
			t.Errorf("Expected to write %d bytes, wrote %d", len(message), n)
		}

		// Read response
		buf := make([]byte, 1024)
		n, err = transport.Read(buf)
		if err != nil {
			t.Fatalf("Read failed: %v", err)
		}

		response := string(buf[:n])
		if response != message {
			t.Errorf("Expected response %q, got %q", message, response)
		}
	})

	t.Run("invalid command", func(t *testing.T) {
		_, err := NewStdioTransport("nonexistent", nil, false)
		if err == nil {
			t.Error("Expected error for nonexistent command")
		}
	})

	t.Run("command arguments", func(t *testing.T) {
		// Use echo to test argument passing
		transport, err := NewStdioTransport("echo", []string{"test message"}, false)
		if err != nil {
			t.Fatalf("Failed to create transport: %v", err)
		}

		if err := transport.Connect(ctx); err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}

		// Read the echoed message
		buf := make([]byte, 1024)
		n, err := transport.Read(buf)
		if err != nil && err != io.EOF {
			t.Fatalf("Read failed: %v", err)
		}

		response := strings.TrimSpace(string(buf[:n]))
		if response != "test message" {
			t.Errorf("Expected response %q, got %q", "test message", response)
		}
	})

	t.Run("debug logging", func(t *testing.T) {
		transport, err := NewStdioTransport("echo", []string{"debug test"}, true)
		if err != nil {
			t.Fatalf("Failed to create transport: %v", err)
		}

		if !transport.debug {
			t.Error("Debug flag not set")
		}

		if err := transport.Connect(ctx); err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
	})
}
