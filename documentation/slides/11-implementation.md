# Technical Implementation

## MCP Client Implementation

```go
// STDIO Transport
cmd := exec.Command(serverConfig.Command[0], serverConfig.Command[1:]...)
transport := mcp.NewCommandTransport(cmd)

// HTTP Transport
baseURL := fmt.Sprintf("http://%s:%d%s", host, port, path)
transport := mcp.NewStreamableClientTransport(baseURL, &options{})

// SSE Transport
transport := mcp.NewSSEClientTransport(baseURL, &options{})
```

## Provider Interface

```go
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
```

---
