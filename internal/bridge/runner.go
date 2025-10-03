package bridge

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
)

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
				b.Log("Reached EOF, stopping")
				return
			}
			if err != nil {
				errChan <- fmt.Errorf("error reading from stdin: %v", err)
				return
			}

			if err := b.StreamToServer(buffer[:n]); err != nil {
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