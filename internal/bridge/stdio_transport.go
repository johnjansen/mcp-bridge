package bridge

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

// StdioTransport implements an MCP transport that connects to a target process via stdio
type StdioTransport struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
	debug  bool

	mu     sync.Mutex
	closed bool
}

// NewStdioTransport creates a new StdioTransport that will execute the given command
// and connect to its stdio for MCP communication
func NewStdioTransport(command string, args []string, debug bool) (*StdioTransport, error) {
	cmd := exec.Command(command, args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		stdin.Close()
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		stdin.Close()
		stdout.Close()
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	return &StdioTransport{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
		debug:  debug,
	}, nil
}

// Connect starts the target process and establishes stdio connections
func (t *StdioTransport) Connect(ctx context.Context) error {
	if err := t.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start target process: %w", err)
	}

	// Start stderr logging goroutine if debug enabled
	if t.debug {
		go func() {
			buf := make([]byte, 4096)
			for {
				n, err := t.stderr.Read(buf)
				if err != nil {
					if err != io.EOF {
						fmt.Printf("Target stderr error: %v\n", err)
					}
					return
				}
				if n > 0 {
					fmt.Printf("Target stderr: %s", buf[:n])
				}
			}
		}()
	}

	return nil
}

// Read implements io.Reader for the transport
func (t *StdioTransport) Read(p []byte) (n int, err error) {
	return t.stdout.Read(p)
}

// Write implements io.Writer for the transport
func (t *StdioTransport) Write(p []byte) (n int, err error) {
	t.mu.Lock()
	if t.closed {
		t.mu.Unlock()
		return 0, io.ErrClosedPipe
	}
	t.mu.Unlock()
	return t.stdin.Write(p)
}

// Close implements io.Closer for the transport
func (t *StdioTransport) Close() error {
	t.mu.Lock()
	if t.closed {
		t.mu.Unlock()
		return nil
	}
	t.closed = true
	t.mu.Unlock()

	// Close pipes and wait for process to exit
	t.stdin.Close()
	t.stdout.Close()
	t.stderr.Close()
	return t.cmd.Wait()
}
