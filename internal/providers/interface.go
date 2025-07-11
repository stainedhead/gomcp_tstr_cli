package providers

import (
	"context"
	"io"
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Tool represents an available tool
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// ToolCall represents a tool call request
type ToolCall struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// ToolResult represents the result of a tool call
type ToolResult struct {
	ID      string      `json:"id"`
	Content interface{} `json:"content"`
	Error   string      `json:"error,omitempty"`
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Messages    []Message    `json:"messages"`
	Tools       []Tool       `json:"tools,omitempty"`
	Stream      bool         `json:"stream,omitempty"`
	Temperature *float64     `json:"temperature,omitempty"`
	MaxTokens   *int         `json:"max_tokens,omitempty"`
	SystemPrompt string      `json:"system_prompt,omitempty"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
	Finished  bool       `json:"finished"`
	Error     string     `json:"error,omitempty"`
}

// Provider defines the interface for AI model providers
type Provider interface {
	// Name returns the provider name
	Name() string

	// Chat sends a chat request and returns a response
	Chat(ctx context.Context, request *ChatRequest) (*ChatResponse, error)

	// ChatStream sends a chat request and streams the response
	ChatStream(ctx context.Context, request *ChatRequest) (<-chan *ChatResponse, error)

	// ValidateConfig validates the provider configuration
	ValidateConfig() error

	// Close closes any resources used by the provider
	Close() error
}

// StreamWriter is a helper interface for writing streaming responses
type StreamWriter interface {
	io.Writer
	Flush() error
}
