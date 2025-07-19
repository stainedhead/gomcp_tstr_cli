# Extensibility

## Adding New Providers

```go
type NewProvider struct {
    // Provider-specific fields
}

func (p *NewProvider) Name() string {
    return "new_provider"
}

// Implement other interface methods...

// Register in factory
func init() {
    RegisterProviderFactory("new_provider", func(config interface{}) (Provider, error) {
        // Create and return new provider instance
    })
}
```

## Adding New Transport Types

1. Extend the `MCPTransport` configuration
2. Add transport creation logic in `initializeServer`
3. Implement the transport using the MCP SDK

---
