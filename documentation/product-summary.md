# mcp_tstr Product Summary

## Overview

`mcp_tstr` is a comprehensive command-line interface tool designed for testing and interacting with Model Context Protocol (MCP) servers. It serves as both a diagnostic tool for MCP server developers and a bridge between AI models and MCP-enabled tools and resources.

## Key Features

### Multi-Protocol Transport Support
- **STDIO Transport**: Direct communication with subprocess-based MCP servers
- **HTTP Transport**: RESTful communication with HTTP-based MCP servers  
- **SSE Transport**: Server-Sent Events for real-time streaming communication
- **Flexible Configuration**: JSON-based server configuration compatible with MCP SDK standards

### Comprehensive Discovery Interface
- **Tool Discovery**: List and inspect available tools from MCP servers
- **Resource Discovery**: Enumerate accessible resources and their metadata
- **Prompt Discovery**: View available prompts and their parameters
- **Unified Listing**: Combined view of all server capabilities

### Interactive Tool Execution
- **Direct Tool Calls**: Execute MCP tools with JSON-formatted parameters
- **Parameter Validation**: Automatic validation of tool parameters
- **Result Formatting**: Pretty-printed JSON output with optional raw mode
- **Error Handling**: Comprehensive error reporting and debugging information

### AI Model Integration
- **Provider Abstraction**: Pluggable architecture for multiple AI providers
- **Streaming Support**: Real-time response streaming for better user experience
- **Tool Integration**: Automatic tool availability injection into AI conversations
- **Multi-Server Chat**: Simultaneous access to tools from multiple MCP servers

### Configuration Management
- **YAML Configuration**: Human-readable configuration files (mcp_tstr.config)
- **Environment Variables**: Override configuration with environment variables
- **Default Values**: Sensible defaults for quick setup
- **Provider-Specific Settings**: Tailored configuration for each AI provider

## Technical Architecture

### Core Components

#### MCP Client Manager
- Manages connections to multiple MCP servers simultaneously
- Handles different transport protocols transparently
- Provides connection pooling and error recovery
- Implements timeout and retry logic

#### Provider Interface
- Abstract interface for AI model providers
- Supports both streaming and non-streaming responses
- Handles tool calling and function execution
- Extensible architecture for adding new providers

#### Chat Session Management
- Maintains conversation context and history
- Manages tool availability and execution
- Provides interactive user interface
- Handles streaming responses and tool integration

#### Configuration System
- Hierarchical configuration with defaults and overrides
- Environment variable integration
- Validation and error reporting
- Hot-reload capabilities for development

### Supported AI Providers

#### Ollama (Implemented)
- Local model execution
- Streaming response support
- Tool calling integration
- Configurable endpoints and models

#### Planned Providers
- **OpenAI**: GPT models with function calling
- **AWS Bedrock**: Claude and other foundation models
- **Google AI**: Gemini models with tool support
- **Anthropic**: Direct Claude API integration

## Use Cases

### MCP Server Development
- **Testing**: Validate server implementations and tool functionality
- **Debugging**: Inspect server responses and error conditions
- **Documentation**: Generate tool and resource documentation
- **Integration Testing**: Verify compatibility with different transports

### AI Application Development
- **Prototyping**: Quickly test AI models with MCP tools
- **Tool Discovery**: Explore available tools and their capabilities
- **Integration Testing**: Validate AI-tool interactions
- **Performance Testing**: Measure response times and throughput

### System Administration
- **Health Monitoring**: Check MCP server availability and status
- **Configuration Validation**: Verify server configurations
- **Troubleshooting**: Diagnose connection and execution issues
- **Performance Monitoring**: Track server response times

## Technical Specifications

### Supported Protocols
- **MCP Version**: Compatible with MCP SDK v0.1.0
- **Transport Protocols**: STDIO, HTTP, SSE
- **Message Format**: JSON-RPC 2.0
- **Authentication**: Provider-specific authentication mechanisms

### Performance Characteristics
- **Concurrent Connections**: Supports multiple simultaneous server connections
- **Streaming**: Real-time response streaming for improved user experience
- **Memory Efficiency**: Minimal memory footprint with connection pooling
- **Error Recovery**: Automatic retry and failover mechanisms

### Security Features
- **Credential Management**: Secure handling of API keys and tokens
- **Environment Isolation**: Separate environments for different server configurations
- **Audit Logging**: Comprehensive logging of all interactions
- **Input Validation**: Strict validation of user inputs and server responses

## Development and Extensibility

### Code Quality
- **Test Coverage**: Comprehensive unit and integration tests
- **Linting**: Strict code quality enforcement with golangci-lint
- **Documentation**: Extensive inline documentation and examples
- **Error Handling**: Robust error handling and user feedback

### Extensibility Points
- **Provider Interface**: Easy addition of new AI providers
- **Transport Layer**: Pluggable transport implementations
- **Configuration**: Extensible configuration schema
- **Middleware**: Request/response middleware support

### Build and Deployment
- **Go Modules**: Modern dependency management
- **Cross-Platform**: Builds on macOS, Linux, and Windows
- **Static Binary**: Single executable with no external dependencies
- **Container Support**: Docker-ready for containerized deployments

## Future Roadmap

### Short Term
- Complete implementation of remaining AI providers
- Enhanced error reporting and debugging tools
- Performance optimizations and connection pooling
- Additional transport protocol support

### Medium Term
- Web-based UI for visual server interaction
- Plugin system for custom tools and providers
- Advanced configuration management
- Metrics and monitoring integration

### Long Term
- Distributed MCP server orchestration
- Advanced AI model routing and load balancing
- Integration with popular development tools
- Enterprise features and security enhancements

## Conclusion

`mcp_tstr` represents a comprehensive solution for MCP server testing and AI model integration. Its flexible architecture, extensive feature set, and focus on developer experience make it an essential tool for anyone working with MCP servers or building AI-powered applications that leverage external tools and resources.

The tool's design prioritizes both ease of use for quick testing scenarios and extensibility for complex integration requirements, making it suitable for individual developers, teams, and enterprise deployments.
