# Configuration Examples

## mcp_tstr.config

```yaml
# Default server to use when none is specified
default_server: "filesystem"

# Default provider to use for chat functionality
default_provider: "ollama"

# Default model to use
default_model: "llama2"

# Logging configuration
logging:
  level: "info"
  to_file: false

# Provider configurations
providers:
  # Ollama configuration
  ollama:
    endpoint: "http://localhost:11434"
    model: "llama2"
```

## mcp.json

```json
{
  "servers": {
    "filesystem": {
      "name": "filesystem",
      "command": ["mcp-server-filesystem"],
      "args": ["/tmp"],
      "transport": {
        "type": "stdio"
      }
    },
    "web_server": {
      "name": "web_server",
      "transport": {
        "type": "http",
        "host": "localhost",
        "port": 8080,
        "path": "/mcp"
      }
    }
  }
}
```

---
