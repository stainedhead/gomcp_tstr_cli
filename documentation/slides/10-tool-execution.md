# Tool Execution Flow

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│   CLI    │     │  Command │     │   MCP    │     │   MCP    │
│  User    │     │  Handler │     │  Manager │     │  Server  │
└────┬─────┘     └────┬─────┘     └────┬─────┘     └────┬─────┘
     │                │                │                │
     │ call-tool      │                │                │
     │ --name=X       │                │                │
     │ --params={}    │                │                │
     │───────────────►│                │                │
     │                │                │                │
     │                │ GetClient()    │                │
     │                │───────────────►│                │
     │                │                │                │
     │                │                │ CallTool(X, {})│
     │                │                │───────────────►│
     │                │                │                │
     │                │                │ Tool Result    │
     │                │                │◄───────────────│
     │                │                │                │
     │                │ Format Output  │                │
     │                │◄───────────────│                │
     │                │                │                │
     │ Display Result │                │                │
     │◄───────────────│                │                │
     │                │                │                │
└──────────┘     └──────────┘     └──────────┘     └──────────┘
```

## CLI Commands
- `list-tools`: Discover available tools
- `list-resources`: View accessible resources
- `list-prompts`: See available prompts
- `call-tool`: Execute specific tool with parameters
- `chat`: Start interactive AI session with tool access

---
