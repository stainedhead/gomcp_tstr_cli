# Configuration System

## Class Structure

```
┌───────────────────────────┐      ┌───────────────────────────┐
│         Config            │      │        MCPConfig          │
├───────────────────────────┤      ├───────────────────────────┤
│ - DefaultServer: string   │      │ - Servers: map[string]    │
│ - DefaultProvider: string │      │              MCPServer    │
│ - DefaultModel: string    │      └───────────────────────────┘
│ - Providers: map          │                   │
│ - Logging: LoggingConfig  │                   │
├───────────────────────────┤      ┌───────────┴───────────────┐
│ + Load(): (*Config, err)  │      │        MCPServer          │
│ + GetProviderConfig()     │      ├───────────────────────────┤
└───────────────────────────┘      │ - Name: string            │
                                   │ - Command: []string       │
                                   │ - Transport: MCPTransport │
                                   └───────────────────────────┘
```

## Configuration Files
- **mcp_tstr.config**: Main application settings
- **mcp.json**: MCP server definitions

---
