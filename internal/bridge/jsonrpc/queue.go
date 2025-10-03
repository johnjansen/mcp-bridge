package jsonrpc

import (
	"context"
	"fmt"
	"sync"
)

// ResponseHandler receives a response for a request
type ResponseHandler chan Message

// MessageQueue manages pending requests and their response handlers
type MessageQueue struct {
	mu       sync.RWMutex
	pending  map[interface{}]ResponseHandler
	closed   bool
	onClosed chan struct{}
}

// NewMessageQueue creates a new message queue
func NewMessageQueue() *MessageQueue {
	return &MessageQueue{
		pending:  make(map[interface{}]ResponseHandler),
		onClosed: make(chan struct{}),
	}
}

// AddRequest registers a request and returns a channel that will receive the response
func (q *MessageQueue) AddRequest(id interface{}) (ResponseHandler, error) {
	if id == nil {
		return nil, fmt.Errorf("request ID cannot be nil")
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return nil, fmt.Errorf("message queue is closed")
	}

	if _, exists := q.pending[id]; exists {
		return nil, fmt.Errorf("duplicate request ID: %v", id)
	}

	// Create buffered channel to avoid blocking
	handler := make(ResponseHandler, 1)
	q.pending[id] = handler
	return handler, nil
}

// HandleResponse delivers a response to the waiting handler
func (q *MessageQueue) HandleResponse(msg Message) error {
	id := msg.GetID()
	if id == nil {
		return fmt.Errorf("response ID cannot be nil")
	}

	q.mu.Lock()
	handler, exists := q.pending[id]
	if exists {
		delete(q.pending, id)
	}
	q.mu.Unlock()

	if !exists {
		return fmt.Errorf("no handler for response ID: %v", id)
	}

	// Try to send the response, but don't block if the handler is gone
	select {
	case handler <- msg:
		return nil
	default:
		return fmt.Errorf("response handler unavailable for ID: %v", id)
	}
}

// Close closes all pending handlers and prevents new requests
func (q *MessageQueue) Close() {
	q.mu.Lock()
	if !q.closed {
		q.closed = true
		close(q.onClosed)

		// Close all pending handlers
		for id, handler := range q.pending {
			close(handler)
			delete(q.pending, id)
		}
	}
	q.mu.Unlock()
}

// IsClosed returns true if the queue has been closed
func (q *MessageQueue) IsClosed() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.closed
}

// WaitResponse waits for a response to a specific request
func (q *MessageQueue) WaitResponse(ctx context.Context, id interface{}) (Message, error) {
	handler, err := q.AddRequest(id)
	if err != nil {
		return nil, fmt.Errorf("failed to register request: %w", err)
	}

	select {
	case <-ctx.Done():
		q.mu.Lock()
		delete(q.pending, id)
		q.mu.Unlock()
		return nil, ctx.Err()

	case <-q.onClosed:
		return nil, fmt.Errorf("message queue closed")

	case msg, ok := <-handler:
		if !ok {
			return nil, fmt.Errorf("response handler closed")
		}
		return msg, nil
	}
}

// PendingCount returns the number of pending requests
func (q *MessageQueue) PendingCount() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.pending)
}

// CancelRequest cancels a pending request by ID
func (q *MessageQueue) CancelRequest(id interface{}) {
	q.mu.Lock()
	if handler, exists := q.pending[id]; exists {
		close(handler)
		delete(q.pending, id)
	}
	q.mu.Unlock()
}
