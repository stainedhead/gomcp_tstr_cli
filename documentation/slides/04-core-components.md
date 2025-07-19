# Core Components

## Four Key Components:

1. **MCP Client Manager**
   - Manages connections to MCP servers
   - Supports multiple transport protocols
   - Handles connection lifecycle

2. **Provider Interface**
   - Abstracts AI model providers
   - Standardizes chat and tool interactions
   - Currently implements Ollama provider

3. **Chat Session Manager**
   - Orchestrates user, AI, and tool interactions
   - Manages conversation context
   - Handles tool execution

4. **Configuration System**
   - Manages application settings
   - Handles server and provider configurations
   - Supports environment variables

---
