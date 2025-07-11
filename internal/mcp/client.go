package mcp

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sirupsen/logrus"

	"mcp_tstr/internal/config"
)

// Client represents an MCP client wrapper
type Client struct {
	name    string
	config  config.MCPServer
	client  *mcp.Client
	session *mcp.ClientSession
	logger  *logrus.Entry
}

// Manager manages multiple MCP clients
type Manager struct {
	clients map[string]*Client
	logger  *logrus.Logger
}

// NewManager creates a new MCP client manager
func NewManager(logger *logrus.Logger) *Manager {
	return &Manager{
		clients: make(map[string]*Client),
		logger:  logger,
	}
}

// InitializeServers initializes MCP servers based on configuration
func (m *Manager) InitializeServers(mcpConfig *config.MCPConfig, serverNames []string) error {
	if len(serverNames) == 0 {
		// Initialize all servers
		for name := range mcpConfig.Servers {
			serverNames = append(serverNames, name)
		}
	}

	if len(serverNames) == 0 {
		return fmt.Errorf("no MCP servers available for use")
	}

	for _, name := range serverNames {
		serverConfig, exists := mcpConfig.Servers[name]
		if !exists {
			return fmt.Errorf("server %s not found in configuration", name)
		}

		client, err := m.initializeServer(name, serverConfig)
		if err != nil {
			m.logger.WithError(err).Errorf("Failed to initialize server %s", name)
			continue
		}

		m.clients[name] = client
		m.logger.Infof("Successfully initialized MCP server: %s", name)
	}

	if len(m.clients) == 0 {
		return fmt.Errorf("no MCP servers could be initialized")
	}

	return nil
}

// initializeServer initializes a single MCP server
func (m *Manager) initializeServer(name string, serverConfig config.MCPServer) (*Client, error) {
	logger := m.logger.WithField("server", name)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var transport mcp.Transport
	var err error

	switch serverConfig.Transport.Type {
	case "stdio":
		transport, err = m.createStdioTransport(serverConfig)
	case "http", "sse":
		transport, err = m.createHTTPTransport(serverConfig)
	default:
		return nil, fmt.Errorf("unsupported transport type: %s", serverConfig.Transport.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create %s transport: %w", serverConfig.Transport.Type, err)
	}

	// Create MCP client
	mcpClient := mcp.NewClient("mcp_tstr", "1.0.0", &mcp.ClientOptions{})

	// Connect to the server
	session, err := mcpClient.Connect(ctx, transport)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	client := &Client{
		name:    name,
		config:  serverConfig,
		client:  mcpClient,
		session: session,
		logger:  logger,
	}

	// Test connection with ping
	if err := client.Ping(ctx); err != nil {
		logger.WithError(err).Warn("Server ping failed, but continuing")
	}

	return client, nil
}

// createStdioTransport creates a STDIO transport
func (m *Manager) createStdioTransport(serverConfig config.MCPServer) (mcp.Transport, error) {
	if len(serverConfig.Command) == 0 {
		return nil, fmt.Errorf("command is required for stdio transport")
	}

	// Create command
	cmd := exec.Command(serverConfig.Command[0], serverConfig.Command[1:]...)
	
	// Set environment variables if provided
	if len(serverConfig.Env) > 0 {
		env := cmd.Environ()
		for key, value := range serverConfig.Env {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
		cmd.Env = env
	}

	return mcp.NewCommandTransport(cmd), nil
}

// createHTTPTransport creates an HTTP or SSE transport
func (m *Manager) createHTTPTransport(serverConfig config.MCPServer) (mcp.Transport, error) {
	if serverConfig.Transport.Host == "" {
		return nil, fmt.Errorf("host is required for http/sse transport")
	}

	baseURL := fmt.Sprintf("http://%s:%d%s", 
		serverConfig.Transport.Host, 
		serverConfig.Transport.Port, 
		serverConfig.Transport.Path)

	if serverConfig.Transport.Type == "sse" {
		return mcp.NewSSEClientTransport(baseURL, &mcp.SSEClientTransportOptions{}), nil
	}

	// For HTTP, we'll use the streamable client transport
	return mcp.NewStreamableClientTransport(baseURL, &mcp.StreamableClientTransportOptions{}), nil
}

// GetClient returns a client by name
func (m *Manager) GetClient(name string) (*Client, error) {
	client, exists := m.clients[name]
	if !exists {
		return nil, fmt.Errorf("client %s not found", name)
	}
	return client, nil
}

// GetAllClients returns all initialized clients
func (m *Manager) GetAllClients() map[string]*Client {
	return m.clients
}

// Close closes all MCP clients
func (m *Manager) Close() error {
	var lastErr error
	for name, client := range m.clients {
		if err := client.Close(); err != nil {
			m.logger.WithError(err).Errorf("Failed to close client %s", name)
			lastErr = err
		}
	}
	return lastErr
}

// Ping sends a ping request to the server
func (c *Client) Ping(ctx context.Context) error {
	return c.session.Ping(ctx, &mcp.PingParams{})
}

// ListTools returns the tools available on this server
func (c *Client) ListTools(ctx context.Context) (*mcp.ListToolsResult, error) {
	return c.session.ListTools(ctx, &mcp.ListToolsParams{})
}

// ListResources returns the resources available on this server
func (c *Client) ListResources(ctx context.Context) (*mcp.ListResourcesResult, error) {
	return c.session.ListResources(ctx, &mcp.ListResourcesParams{})
}

// ListPrompts returns the prompts available on this server
func (c *Client) ListPrompts(ctx context.Context) (*mcp.ListPromptsResult, error) {
	return c.session.ListPrompts(ctx, &mcp.ListPromptsParams{})
}

// CallTool executes a tool with the given parameters
func (c *Client) CallTool(ctx context.Context, name string, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	return c.session.CallTool(ctx, &mcp.CallToolParams{
		Name:      name,
		Arguments: arguments,
	})
}

// Close closes the MCP client connection
func (c *Client) Close() error {
	if c.session != nil {
		return c.session.Close()
	}
	return nil
}

// GetName returns the client name
func (c *Client) GetName() string {
	return c.name
}

// GetConfig returns the client configuration
func (c *Client) GetConfig() config.MCPServer {
	return c.config
}
