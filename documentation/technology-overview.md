# mcp_tstr Technology Overview

## Architecture Overview

`mcp_tstr` is a Go-based CLI application designed to test and interact with Model Context Protocol (MCP) servers. The application follows a modular architecture with clear separation of concerns between components. This document provides a technical overview of the system architecture, key components, and their interactions.

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

## Core Components

### 1. MCP Client Manager

The MCP Client Manager is responsible for creating, managing, and communicating with MCP servers through various transport protocols.

#### Class Diagram

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
                            │ + GetName()       │
                            │ + GetConfig()     │
                            └───────────────────┘
```

#### Key Features

- **Multi-Server Support**: Manages connections to multiple MCP servers simultaneously
- **Transport Abstraction**: Supports STDIO, HTTP, and SSE transports
- **Connection Management**: Handles initialization, monitoring, and cleanup of connections
- **Error Handling**: Provides robust error handling and recovery mechanisms

#### Transport Types

1. **STDIO Transport**
   - Uses stdin/stdout for communication with subprocess-based MCP servers
   - Launches and manages server processes
   - Handles environment variables and command-line arguments

2. **HTTP Transport**
   - RESTful communication with HTTP-based MCP servers
   - Supports both synchronous and streaming requests
   - Handles connection pooling and timeout management

3. **SSE Transport**
   - Server-Sent Events for real-time streaming communication
   - Maintains persistent connections for event streaming
   - Handles reconnection and backoff strategies

### 2. Provider Interface

The Provider Interface defines a common abstraction for interacting with different AI model providers.

#### Class Diagram

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

#### Key Features

- **Provider Abstraction**: Common interface for all AI model providers
- **Streaming Support**: Both synchronous and streaming response modes
- **Tool Integration**: Standardized tool calling interface
- **Error Handling**: Consistent error reporting across providers

#### Message and Tool Types

```
┌───────────────────┐      ┌───────────────────┐
│      Message      │      │       Tool        │
├───────────────────┤      ├───────────────────┤
│ - Role: string    │      │ - Name: string    │
│ - Content: string │      │ - Description: str│
└───────────────────┘      │ - Parameters: map │
                           └───────────────────┘
                                     ▲
                                     │
                           ┌─────────┴─────────┐
                           │     ToolCall      │
                           ├───────────────────┤
                           │ - ID: string      │
                           │ - Name: string    │
                           │ - Arguments: map  │
                           └───────────────────┘
```

#### Implemented Providers

1. **Ollama Provider**
   - Local model execution via Ollama API
   - Streaming response support
   - Tool calling integration (when supported by model)
   - Configurable endpoints and models

### 3. Chat Session Manager

The Chat Session Manager orchestrates interactions between users, AI models, and MCP tools.

#### Class Diagram

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

#### Key Features

- **Conversation Management**: Maintains conversation history and context
- **Tool Integration**: Dynamically loads and provides tools to AI models
- **Interactive Interface**: Provides a user-friendly CLI interface
- **Streaming Responses**: Real-time streaming of AI responses

#### Session Flow

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│  User    │     │  Chat    │     │   AI     │     │   MCP    │
│ Input    │────►│ Session  │────►│ Provider │     │ Manager  │
└──────────┘     └────┬─────┘     └────┬─────┘     └────┬─────┘
                      │                │                │
                      │                │                │
                      │   Request      │                │
                      │───────────────►│                │
                      │                │                │
                      │   Response     │                │
                      │◄───────────────│                │
                      │                │                │
                      │          Tool Call              │
                      │───────────────────────────────►│
                      │                │                │
                      │          Tool Result            │
                      │◄───────────────────────────────│
                      │                │                │
┌──────────┐     ┌────┴─────┐     ┌────┴─────┐     ┌────┴─────┐
│  User    │     │  Chat    │     │   AI     │     │   MCP    │
│ Display  │◄────│ Session  │     │ Provider │     │ Manager  │
└──────────┘     └──────────┘     └──────────┘     └──────────┘
```

### 4. Configuration System

The Configuration System manages application settings, server configurations, and provider options.

#### Class Diagram

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
                                   │ - Args: []string          │
                                   │ - Env: map[string]string  │
                                   │ - Transport: MCPTransport │
                                   │ - Extra: map              │
                                   └───────────────────────────┘
                                                │
                                   ┌───────────┴───────────────┐
                                   │       MCPTransport        │
                                   ├───────────────────────────┤
                                   │ - Type: string            │
                                   │ - Host: string            │
                                   │ - Port: int               │
                                   │ - Path: string            │
                                   └───────────────────────────┘
```

#### Key Features

- **YAML Configuration**: Human-readable configuration files
- **Environment Variables**: Support for environment variable overrides
- **Default Values**: Sensible defaults for quick setup
- **Validation**: Configuration validation and error reporting

#### Configuration Files

1. **mcp_tstr.config**
   - Main application configuration
   - Provider settings
   - Logging configuration
   - Default values

2. **mcp.json**
   - MCP server definitions
   - Transport configurations
   - Server command and environment settings

## Command Layer

The Command Layer provides the CLI interface for interacting with the application.

### Command Structure

```
mcp_tstr
  ├── list-tools
  ├── list-resources
  ├── list-prompts
  ├── list-all
  ├── list-servers
  ├── call-tool
  ├── ping
  ├── chat
  └── show-mcp-json
```

### Sequence Diagrams

#### Tool Discovery Flow

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│   CLI    │     │  Command │     │   MCP    │     │   MCP    │
│  User    │     │  Handler │     │  Manager │     │  Server  │
└────┬─────┘     └────┬─────┘     └────┬─────┘     └────┬─────┘
     │                │                │                │
     │ list-tools     │                │                │
     │───────────────►│                │                │
     │                │                │                │
     │                │ GetClient()    │                │
     │                │───────────────►│                │
     │                │                │                │
     │                │                │ ListTools()    │
     │                │                │───────────────►│
     │                │                │                │
     │                │                │ Tools Response │
     │                │                │◄───────────────│
     │                │                │                │
     │                │ Format Output  │                │
     │                │◄───────────────│                │
     │                │                │                │
     │ Display Tools  │                │                │
     │◄───────────────│                │                │
     │                │                │                │
┌────┴─────┐     ┌────┴─────┐     ┌────┴─────┐     ┌────┴─────┐
│   CLI    │     │  Command │     │   MCP    │     │   MCP    │
│  User    │     │  Handler │     │  Manager │     │  Server  │
└──────────┘     └──────────┘     └──────────┘     └──────────┘
```

#### Tool Execution Flow

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
┌────┴─────┐     ┌────┴─────┐     ┌────┴─────┐     ┌────┴─────┐
│   CLI    │     │  Command │     │   MCP    │     │   MCP    │
│  User    │     │  Handler │     │  Manager │     │  Server  │
└──────────┘     └──────────┘     └──────────┘     └──────────┘
```

#### Chat Session Flow

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
     │                │                │ Start()        │                │
     │                │                │────────────────┐                │
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
┌────┴─────┐     ┌────┴─────┐     ┌────┴─────┐     ┌────┴─────┐     ┌────┴─────┐
│   CLI    │     │  Command │     │   Chat   │     │   AI     │     │   MCP    │
│  User    │     │  Handler │     │  Session │     │ Provider │     │  Server  │
└──────────┘     └──────────┘     └──────────┘     └──────────┘     └──────────┘
```

## Technical Implementation Details

### MCP Client Implementation

The MCP Client Manager uses the Model Context Protocol Go SDK to communicate with MCP servers. It supports three transport types:

1. **STDIO Transport**:
   ```go
   cmd := exec.Command(serverConfig.Command[0], serverConfig.Command[1:]...)
   transport := mcp.NewCommandTransport(cmd)
   ```

2. **HTTP Transport**:
   ```go
   baseURL := fmt.Sprintf("http://%s:%d%s", serverConfig.Transport.Host, serverConfig.Transport.Port, serverConfig.Transport.Path)
   transport := mcp.NewStreamableClientTransport(baseURL, &mcp.StreamableClientTransportOptions{})
   ```

3. **SSE Transport**:
   ```go
   baseURL := fmt.Sprintf("http://%s:%d%s", serverConfig.Transport.Host, serverConfig.Transport.Port, serverConfig.Transport.Path)
   transport := mcp.NewSSEClientTransport(baseURL, &mcp.SSEClientTransportOptions{})
   ```

### Provider Interface

The Provider Interface defines a common abstraction for AI model providers:

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

### Ollama Provider Implementation

The Ollama Provider implements the Provider interface for local Ollama models:

```go
func (p *OllamaProvider) ChatStream(ctx context.Context, request *ChatRequest) (<-chan *ChatResponse, error) {
    // Convert to Ollama-specific request format
    ollamaReq := p.convertRequest(request)
    ollamaReq.Stream = true

    // Create HTTP request
    reqBody, _ := json.Marshal(ollamaReq)
    httpReq, _ := http.NewRequestWithContext(ctx, "POST", p.endpoint+"/api/chat", bytes.NewReader(reqBody))
    httpReq.Header.Set("Content-Type", "application/json")

    // Send request
    resp, err := p.client.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }

    // Process streaming response
    responseChan := make(chan *ChatResponse, 10)
    go func() {
        defer close(responseChan)
        defer resp.Body.Close()

        decoder := json.NewDecoder(resp.Body)
        for {
            var ollamaResp OllamaResponse
            if err := decoder.Decode(&ollamaResp); err != nil {
                if err == io.EOF {
                    break
                }
                responseChan <- &ChatResponse{Error: fmt.Sprintf("decode error: %v", err)}
                return
            }

            responseChan <- &ChatResponse{
                Content:  ollamaResp.Message.Content,
                Finished: ollamaResp.Done,
            }

            if ollamaResp.Done {
                break
            }
        }
    }()

    return responseChan, nil
}
```

### Chat Session Implementation

The Chat Session manages the interaction between the user, AI provider, and MCP tools:

```go
func (s *Session) processMessage(ctx context.Context) error {
    // Create chat request with tools
    request := &providers.ChatRequest{
        Messages:     s.messages,
        Tools:        s.tools,
        Stream:       true,
        SystemPrompt: s.systemPrompt,
    }

    // Stream response from AI provider
    responseChan, err := s.provider.ChatStream(ctx, request)
    if err != nil {
        return fmt.Errorf("failed to start chat stream: %w", err)
    }

    // Process streaming response
    fmt.Print("Assistant: ")
    var fullResponse strings.Builder

    for response := range responseChan {
        if response.Error != "" {
            return fmt.Errorf("chat error: %s", response.Error)
        }

        // Print streaming content
        fmt.Print(response.Content)
        fullResponse.WriteString(response.Content)

        // Handle tool calls if present
        if len(response.ToolCalls) > 0 {
            for _, toolCall := range response.ToolCalls {
                if err := s.handleToolCall(ctx, toolCall); err != nil {
                    fmt.Printf("\n[Tool call failed: %v]", err)
                }
            }
        }

        if response.Finished {
            break
        }
    }

    // Add response to conversation history
    s.messages = append(s.messages, providers.Message{
        Role:    "assistant",
        Content: fullResponse.String(),
    })

    return nil
}
```

## Configuration Format

### mcp_tstr.config

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

### mcp.json

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
    }
  }
}
```

## Performance Considerations

### Connection Management

- **Connection Pooling**: Reuses connections to MCP servers when possible
- **Timeout Handling**: Implements configurable timeouts for all operations
- **Error Recovery**: Automatic retry and reconnection logic
- **Resource Cleanup**: Proper cleanup of resources on shutdown

### Memory Management

- **Streaming Responses**: Uses Go channels for efficient streaming
- **Buffer Management**: Careful buffer sizing to prevent memory issues
- **Resource Limits**: Configurable limits on concurrent connections

### Concurrency

- **Context-Based Cancellation**: All operations support context cancellation
- **Goroutine Management**: Careful management of goroutines to prevent leaks
- **Channel Sizing**: Appropriate buffer sizes for communication channels

## Security Considerations

### API Key Management

- **Environment Variables**: Support for environment variable-based API keys
- **Configuration Isolation**: Local configuration files not tracked in git
- **Secure Defaults**: Sensible security defaults

### Input Validation

- **Parameter Validation**: Strict validation of user inputs
- **JSON Schema Validation**: Validation of tool parameters against schemas
- **Error Handling**: Secure error handling that doesn't leak sensitive information

## Testing Strategy

### Unit Tests

- **Component Testing**: Tests for individual components
- **Mock Interfaces**: Mock implementations of interfaces for testing
- **Error Handling**: Tests for error conditions and edge cases

### Integration Tests

- **End-to-End Testing**: Tests for complete workflows
- **Server Mocking**: Mock MCP servers for testing
- **Provider Mocking**: Mock AI providers for testing

## Extensibility

### Adding New Providers

To add a new AI provider:

1. Implement the `Provider` interface
2. Add provider-specific configuration
3. Register the provider in the factory

Example:

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

### Adding New Transport Types

To add a new transport type:

1. Extend the `MCPTransport` configuration
2. Add transport creation logic in `initializeServer`
3. Implement the transport using the MCP SDK

## Deployment Considerations

### Binary Distribution

- **Single Binary**: Compiled as a single statically-linked binary
- **Cross-Platform**: Builds for multiple platforms (Linux, macOS, Windows)
- **Minimal Dependencies**: No external runtime dependencies

### Configuration Management

- **Example Configurations**: Provided example configurations
- **Environment Variables**: Support for environment-based configuration
- **Validation**: Configuration validation on startup

## Future Technical Enhancements

### Planned Technical Improvements

1. **Connection Pooling**: Enhanced connection pooling for better performance
2. **Metrics Collection**: Instrumentation for performance monitoring
3. **Plugin System**: Plugin architecture for custom extensions
4. **Web Interface**: Optional web-based UI for visual interaction
5. **Authentication**: Enhanced authentication mechanisms for MCP servers
6. **Distributed Operation**: Support for distributed MCP server orchestration

## Conclusion

The `mcp_tstr` tool provides a comprehensive solution for testing and interacting with MCP servers. Its modular architecture, extensive feature set, and focus on developer experience make it an essential tool for MCP server development and AI model integration.

The tool's design prioritizes:

- **Flexibility**: Support for multiple transport types and AI providers
- **Extensibility**: Easy addition of new providers and features
- **Robustness**: Comprehensive error handling and recovery
- **Performance**: Efficient resource usage and streaming support
- **Security**: Secure handling of credentials and configurations

This technical overview provides a foundation for understanding the architecture and implementation details of the `mcp_tstr` tool, enabling developers to effectively use, extend, and contribute to the project.
