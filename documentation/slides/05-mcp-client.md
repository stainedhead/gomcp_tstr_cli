# MCP Client Manager

## Class Structure

```
┌───────────────────┐       ┌───────────────────┐
│     Manager       │       │      Client       │
├───────────────────┤       ├───────────────────┤
│ - clients: map    │       │ - name: string    │
│ - logger: Logger  │       │ - config: Config  │
├───────────────────┤       │ - client: *Client │
│ + NewManager()    │       │ - session: *Session│
│ + InitializeServers()◄────┤ - logger: Logger  │
│ + GetClient()     │       ├───────────────────┤
│ + GetAllClients() │       │ + Ping()          │
│ + Close()         │       │ + ListTools()     │
└───────────────────┘       │ + ListResources() │
                            │ + ListPrompts()   │
                            │ + CallTool()      │
                            │ + Close()         │
                            └───────────────────┘
```

## Transport Types
- **STDIO**: For subprocess-based servers
- **HTTP**: For RESTful API servers
- **SSE**: For event streaming servers

---
