package jsonrpc

import (
	"encoding/json"
	"fmt"
)

// Version is the JSON-RPC version string
const Version = "2.0"

// Message represents a JSON-RPC 2.0 message
type Message interface {
	// GetID returns the message ID. May be nil for notifications.
	GetID() interface{}
	// GetVersion returns the JSON-RPC version string
	GetVersion() string
	// GetType returns the type of message (request, response, notification, error)
	GetType() MessageType
}

// MessageType indicates the type of JSON-RPC message
type MessageType int

const (
	RequestType MessageType = iota
	ResponseType
	NotificationType
	ErrorType
)

// Request represents a JSON-RPC 2.0 request
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      interface{}     `json:"id"`
}

func (r *Request) GetID() interface{}   { return r.ID }
func (r *Request) GetVersion() string   { return r.JSONRPC }
func (r *Request) GetType() MessageType { return RequestType }

// Response represents a JSON-RPC 2.0 response
type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *Error          `json:"error,omitempty"`
	ID      interface{}     `json:"id"`
}

func (r *Response) GetID() interface{}   { return r.ID }
func (r *Response) GetVersion() string   { return r.JSONRPC }
func (r *Response) GetType() MessageType { return ResponseType }

// Error represents a JSON-RPC 2.0 error object
type Error struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// Standard error codes from JSON-RPC 2.0 spec
const (
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
)

// Create an error response with the given code and message
func NewError(id interface{}, code int, msg string, data interface{}) *Response {
	var rawData json.RawMessage
	if data != nil {
		if raw, err := json.Marshal(data); err == nil {
			rawData = raw
		}
	}

	return &Response{
		JSONRPC: Version,
		ID:      id,
		Error: &Error{
			Code:    code,
			Message: msg,
			Data:    rawData,
		},
	}
}

// Parse parses a JSON-RPC message and returns the appropriate type
func Parse(data []byte) (Message, error) {
	// First try to parse as a general JSON object to determine type
	var msg struct {
		JSONRPC string          `json:"jsonrpc"`
		Method  string          `json:"method"`
		Result  json.RawMessage `json:"result"`
		Error   *Error          `json:"error"`
		ID      interface{}     `json:"id"`
	}

	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("failed to parse JSON-RPC message: %w", err)
	}

	// Verify protocol version
	if msg.JSONRPC != Version {
		return nil, fmt.Errorf("unsupported JSON-RPC version: %s", msg.JSONRPC)
	}

	// Determine message type and parse accordingly
	switch {
	case msg.Method != "":
		// It's a request or notification
		var req Request
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, fmt.Errorf("failed to parse request: %w", err)
		}
		return &req, nil

	case msg.Error != nil:
		// It's an error response
		var resp Response
		if err := json.Unmarshal(data, &resp); err != nil {
			return nil, fmt.Errorf("failed to parse error response: %w", err)
		}
		return &resp, nil

	default:
		// It's a success response
		var resp Response
		if err := json.Unmarshal(data, &resp); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		return &resp, nil
	}
}
