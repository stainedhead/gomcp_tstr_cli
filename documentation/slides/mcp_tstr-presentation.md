---
marp: true
theme: default
paginate: true
header: "mcp_tstr - Model Context Protocol Testing Tool"
footer: "© 2025 - mcp_tstr Project"
---

# mcp_tstr

## Model Context Protocol Testing Tool

A comprehensive CLI tool for testing and interacting with MCP servers

---
# Overview

- **Purpose**: Test and interact with Model Context Protocol (MCP) servers
- **Key Features**:
  - Multi-protocol transport support (STDIO, HTTP, SSE)
  - Comprehensive discovery interface
  - Interactive tool execution
  - AI model integration
  - Flexible configuration

---
# System Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                           mcp_tstr CLI                              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌───────────────┐    ┌───────────────┐    ┌───────────────────┐    │
│  │ Command Layer │    │ Configuration │    │ Logging & Utility │    │
│  └───────┬───────┘    └───────┬───────┘    └───────────────────┘    │
│          │                    │                                     │
│          ▼                    ▼                                     │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │                     Core Components                         │    │
│  │                                                             │    │
│  │  ┌───────────────┐    ┌───────────────┐    ┌─────────────┐  │    │
│  │  │  MCP Client   │◄──►│ Chat Session  │◄──►│ AI Provider │  │    │
│  │  │   Manager     │    │   Manager     │    │  Interface  │  │    │
│  │  └───────┬───────┘    └───────────────┘    └──────┬──────┘  │    │
│  │          │                                        │         │    │
│  │          ▼                                        ▼         │    │
│  │  ┌───────────────┐                      ┌─────────────────┐ │    │
│  │  │ Transport     │                      │ Provider        │ │    │
│  │  │ Implementations│                      │ Implementations │ │    │
│  │  └───────────────┘                      └─────────────────┘ │    │
│  │                                                             │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---
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
# Provider Interface

## Class Structure

```
┌───────────────────────────────┐
│        <<interface>>          │
│          Provider             │
├───────────────────────────────┤
│ + Name(): string              │
│ + Chat(): (*ChatResponse, err)│
│ + ChatStream(): (<-chan, err) │
│ + ValidateConfig(): error     │
│ + Close(): error              │
└───────────────┬───────────────┘
                │
                │
┌───────────────┴───────────────┐
│       OllamaProvider          │
├───────────────────────────────┤
│ - endpoint: string            │
│ - model: string               │
│ - client: *http.Client        │
│ - logger: Logger              │
├───────────────────────────────┤
│ + NewOllamaProvider()         │
│ + Name(): string              │
│ + ValidateConfig(): error     │
│ + Chat(): (*ChatResponse, err)│
│ + ChatStream(): (<-chan, err) │
│ + convertRequest()            │
│ + Close(): error              │
└───────────────────────────────┘
```

## Implemented Providers
- **Ollama**: Local model execution
- **Planned**: OpenAI, AWS Bedrock, Google AI, Anthropic

---
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
# Key Workflows

## Chat Session Flow

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│   CLI    │     │  Command │     │   Chat   │     │   AI     │     │   MCP    │
│  User    │     │  Handler │     │  Session │     │ Provider │     │  Server  │
└────┬─────┘     └────┬─────┘     └────┬─────┘     └────┬─────┘     └────┬─────┘
     │                │                │                │                │
     │ chat           │                │                │                │
     │───────────────►│                │                │                │
     │                │                │                │                │
     │                │ NewSession()   │                │                │
     │                │───────────────►│                │                │
     │                │                │                │                │
     │                │                │ LoadTools()    │                │
     │                │                │───────────────────────────────►│
     │                │                │                │                │
     │                │                │ Tools List     │                │
     │                │                │◄───────────────────────────────│
     │                │                │                │                │
     │ User Input     │                │                │                │
     │───────────────────────────────►│                │                │
     │                │                │                │                │
     │                │                │ Chat Request   │                │
     │                │                │───────────────►│                │
     │                │                │                │                │
     │                │                │ AI Response    │                │
     │                │                │◄───────────────│                │
     │                │                │                │                │
     │                │                │ Tool Call      │                │
     │                │                │───────────────────────────────►│
     │                │                │                │                │
     │                │                │ Tool Result    │                │
     │                │                │◄───────────────────────────────│
     │                │                │                │                │
     │ Display Output │                │                │                │
     │◄───────────────────────────────│                │                │
     │                │                │                │                │
└──────────┘     └──────────┘     └──────────┘     └──────────┘     └──────────┘
```

---
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
# Key Considerations

## Performance

- **Connection Pooling**: Reuses connections to MCP servers
- **Streaming Responses**: Uses Go channels for efficient streaming
- **Resource Management**: Proper cleanup of resources
- **Concurrency**: Context-based cancellation for all operations

## Security

- **API Key Management**: Environment variable-based API keys
- **Configuration Isolation**: Local config files not tracked in git
- **Input Validation**: Strict validation of user inputs
- **Error Handling**: Secure error handling

## Testing

- **Unit Tests**: Tests for individual components
- **Integration Tests**: Tests for complete workflows
- **Mock Interfaces**: Mock implementations for testing

---
# Deployment

## Binary Distribution

- **Single Binary**: Compiled as a statically-linked binary
- **Cross-Platform**: Builds for Linux, macOS, Windows
- **Minimal Dependencies**: No external runtime dependencies

## Configuration Management

- **Example Configurations**: Provided in examples/ directory
- **Environment Variables**: Support for environment-based config
- **Validation**: Configuration validation on startup

## Installation

```bash
# Build from source
git clone <repository-url>
cd mcp_tstr
go build -o mcp_tstr .

# Configure
cp examples/mcp_tstr.config mcp_tstr.config
cp examples/mcp.json mcp.json

# Run
./mcp_tstr --help
```

---
# Future Enhancements

## Short Term

- Complete implementation of remaining AI providers
- Enhanced error reporting and debugging tools
- Performance optimizations and connection pooling
- Additional transport protocol support

## Medium Term

- Web-based UI for visual server interaction
- Plugin system for custom tools and providers
- Advanced configuration management
- Metrics and monitoring integration

## Long Term

- Distributed MCP server orchestration
- Advanced AI model routing and load balancing
- Integration with popular development tools
- Enterprise features and security enhancements

---
# Conclusion

## mcp_tstr Key Benefits

- **Comprehensive Testing**: Complete MCP server testing capabilities
- **Flexible Architecture**: Support for multiple transports and providers
- **Developer Experience**: Intuitive CLI interface and rich features
- **Extensibility**: Easy to add new providers and features
- **Production Ready**: Robust error handling and configuration

## Get Involved

- **GitHub**: [github.com/username/mcp_tstr](https://github.com/username/mcp_tstr)
- **Documentation**: See `/documentation` directory
- **Examples**: Check `/examples` directory
- **Issues & PRs**: Contributions welcome!

## Questions?

Thank you!

---
