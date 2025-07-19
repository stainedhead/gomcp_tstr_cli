# Chat Session Manager

## Class Structure

```
┌───────────────────────────────┐
│           Session             │
├───────────────────────────────┤
│ - provider: Provider          │
│ - mcpManager: *Manager        │
│ - messages: []Message         │
│ - tools: []Tool               │
│ - systemPrompt: string        │
│ - logger: Logger              │
├───────────────────────────────┤
│ + NewSession()                │
│ + SetSystemPrompt()           │
│ + LoadTools()                 │
│ + Start()                     │
│ - processMessage()            │
│ - handleToolCall()            │
└───────────────────────────────┘
```

## Key Features
- Conversation history management
- Dynamic tool loading from MCP servers
- Real-time streaming responses
- Tool call handling and execution

---
