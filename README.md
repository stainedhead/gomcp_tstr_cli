# mcp_tstr - MCP Testing CLI Tool

A comprehensive command-line interface tool for testing and interacting with Model Context Protocol (MCP) servers. This tool provides capabilities to discover server features, execute tools, and chat with AI models that have access to MCP server tools and resources.

## Features

- **Multi-Server Support**: Connect to multiple MCP servers simultaneously
- **Transport Flexibility**: Support for STDIO, HTTP, and SSE transports
- **Discovery Interface**: List tools, resources, and prompts from MCP servers
- **Tool Execution**: Execute MCP tools with parameters
- **Interactive Chat**: Chat with AI models that can use MCP tools
- **Provider Support**: Multiple AI model providers (Ollama implemented, others planned)
- **Configuration Management**: YAML configuration with environment variable support
- **Comprehensive Logging**: Configurable logging levels and file output

## Installation

### Prerequisites

- Go 1.23 or later
- Access to MCP servers you want to test

### Build from Source

```bash
git clone <repository-url>
cd mcp_tstr
go build -o mcp_tstr .
```

### Install Dependencies

```bash
go mod tidy
```

## Configuration

**Note**: Configuration files (`mcp_tstr.config` and `mcp.json`) are not tracked in git to prevent accidental commit of sensitive information. Copy the examples to create your local configuration files.

### Application Configuration

Create a `mcp_tstr.config` file in your working directory by copying from the example:

```bash
cp examples/mcp_tstr.config mcp_tstr.config
```

Then edit the configuration file:

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
  
  # Future providers (not yet implemented)
  openai:
    api_key: "${OPENAI_API_KEY}"
    model: "gpt-4"
  
  aws_bedrock:
    region: "us-east-1"
    access_key_id: "${AWS_ACCESS_KEY_ID}"
    secret_access_key: "${AWS_SECRET_ACCESS_KEY}"
    model: "anthropic.claude-3-sonnet-20240229-v1:0"
```

### MCP Server Configuration

Create a `mcp.json` file by copying from the example:

```bash
cp examples/mcp.json mcp.json
```

Then edit the file to define your MCP servers:

```json
{
  "servers": {
    "filesystem": {
      "name": "filesystem",
      "command": ["mcp-server-filesystem"],
      "args": ["/tmp"],
      "transport": {
        "type": "stdio"
      },
      "env": {
        "PATH": "/usr/local/bin:/usr/bin:/bin"
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
    },
    "sse_server": {
      "name": "sse_server",
      "transport": {
        "type": "sse",
        "host": "localhost",
        "port": 9090,
        "path": "/events"
      }
    }
  }
}
```

## Usage

### Global Flags

- `--config`: Specify config file (default: ./mcp_tstr.config)
- `--server, -s`: MCP server name to interact with
- `--provider-name, -p`: Model provider to use for chat
- `--log-level, -l`: Logging level (debug, info, warn, error)
- `--use-all-mcp, -u`: Include all servers in chat session
- `--log-to-file, -f`: Store logs to persistent file
- `--json-raw, -j`: Turn off JSON formatting in results
- `--version, -v`: Show version information
- `--help, -h`: Show help

### Commands

#### Discovery Commands

**List all capabilities:**
```bash
mcp_tstr list-all --server filesystem
```

**List tools:**
```bash
mcp_tstr list-tools --server filesystem
```

**List resources:**
```bash
mcp_tstr list-resources --server filesystem
```

**List prompts:**
```bash
mcp_tstr list-prompts --server filesystem
```

**List configured servers:**
```bash
mcp_tstr list-servers
```

**Show MCP configuration:**
```bash
mcp_tstr show-mcp-json
```

#### Interaction Commands

**Execute a tool:**
```bash
mcp_tstr call-tool --server filesystem --name "read_file" --params '{"path":"/tmp/test.txt"}'
```

**Test server connectivity:**
```bash
mcp_tstr ping --server filesystem
```

**Start interactive chat:**
```bash
mcp_tstr chat --provider-name ollama --server filesystem
```

**Chat with all servers:**
```bash
mcp_tstr chat --provider-name ollama --use-all-mcp
```

### Environment Variables

You can override configuration values using environment variables:

- `OLLAMA_ENDPOINT`: Ollama server endpoint
- `OLLAMA_MODEL`: Ollama model name
- `OPENAI_API_KEY`: OpenAI API key
- `AWS_ACCESS_KEY_ID`: AWS access key
- `AWS_SECRET_ACCESS_KEY`: AWS secret key
- `GOOGLE_AI_API_KEY`: Google AI API key
- `ANTHROPIC_API_KEY`: Anthropic API key

## Transport Types

### STDIO Transport

For MCP servers that communicate via stdin/stdout:

```json
{
  "transport": {
    "type": "stdio"
  },
  "command": ["python", "-m", "my_mcp_server"],
  "env": {
    "PYTHONPATH": "/path/to/server"
  }
}
```

### HTTP Transport

For MCP servers accessible via HTTP:

```json
{
  "transport": {
    "type": "http",
    "host": "localhost",
    "port": 8080,
    "path": "/mcp"
  }
}
```

### SSE Transport

For MCP servers using Server-Sent Events:

```json
{
  "transport": {
    "type": "sse",
    "host": "localhost",
    "port": 9090,
    "path": "/events"
  }
}
```

## AI Model Providers

### Ollama (Implemented)

Supports local Ollama installations:

```yaml
providers:
  ollama:
    endpoint: "http://localhost:11434"
    model: "llama2"
```

### Future Providers

The following providers are planned but not yet implemented:

- **OpenAI**: GPT models with function calling
- **AWS Bedrock**: Claude and other models
- **Google AI**: Gemini models
- **Anthropic**: Direct Claude API

## Chat Features

The interactive chat session provides:

- **Tool Integration**: AI models can automatically use MCP tools
- **Streaming Responses**: Real-time response streaming
- **Multi-Server Support**: Access tools from multiple MCP servers
- **Conversation History**: Maintains context throughout the session
- **Exit Commands**: Type `bye`, `exit`, `end`, or `quit` to end

## Development

### Running Tests

```bash
go test ./...
```

### Linting

```bash
golangci-lint run
```

### Building

```bash
go build -o mcp_tstr .
```

## Troubleshooting

### Common Issues

1. **Server Connection Failed**
   - Verify MCP server is running
   - Check transport configuration
   - Ensure correct command path for STDIO servers

2. **Provider Not Found**
   - Check provider configuration in config file
   - Verify environment variables are set
   - Ensure provider service is accessible

3. **Tool Execution Failed**
   - Verify tool name exists on server
   - Check parameter format (must be valid JSON)
   - Ensure server supports the tool

### Debug Mode

Enable debug logging for detailed information:

```bash
mcp_tstr --log-level debug list-tools --server filesystem
```

### Log to File

Store logs for analysis:

```bash
mcp_tstr --log-to-file chat --provider-name ollama
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run linting and tests
6. Submit a pull request

## License

[Add your license information here]

## Support

For issues and questions:

1. Check the troubleshooting section
2. Review the configuration examples
3. Enable debug logging for detailed error information
4. Open an issue with detailed information about your setup and the problem
