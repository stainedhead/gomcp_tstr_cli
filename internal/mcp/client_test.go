package mcp

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"mcp_tstr/internal/config"
)

func TestNewManager(t *testing.T) {
	logger := logrus.New()
	manager := NewManager(logger)

	assert.NotNil(t, manager)
	assert.NotNil(t, manager.clients)
	assert.Equal(t, logger, manager.logger)
	assert.Empty(t, manager.clients)
}

func TestManagerGetClient(t *testing.T) {
	logger := logrus.New()
	manager := NewManager(logger)

	// Test getting non-existent client
	client, err := manager.GetClient("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "not found")
}

func TestCreateStdioTransport(t *testing.T) {
	logger := logrus.New()
	manager := NewManager(logger)

	tests := []struct {
		name         string
		serverConfig config.MCPServer
		expectError  bool
	}{
		{
			name: "valid stdio config",
			serverConfig: config.MCPServer{
				Command: []string{"echo", "hello"},
				Transport: config.MCPTransport{
					Type: "stdio",
				},
			},
			expectError: false,
		},
		{
			name: "missing command",
			serverConfig: config.MCPServer{
				Command: []string{},
				Transport: config.MCPTransport{
					Type: "stdio",
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport, err := manager.createStdioTransport(tt.serverConfig)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, transport)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, transport)
			}
		})
	}
}

func TestCreateHTTPTransport(t *testing.T) {
	logger := logrus.New()
	manager := NewManager(logger)

	tests := []struct {
		name         string
		serverConfig config.MCPServer
		expectError  bool
	}{
		{
			name: "valid http config",
			serverConfig: config.MCPServer{
				Transport: config.MCPTransport{
					Type: "http",
					Host: "localhost",
					Port: 8080,
					Path: "/mcp",
				},
			},
			expectError: false,
		},
		{
			name: "valid sse config",
			serverConfig: config.MCPServer{
				Transport: config.MCPTransport{
					Type: "sse",
					Host: "localhost",
					Port: 9090,
					Path: "/events",
				},
			},
			expectError: false,
		},
		{
			name: "missing host",
			serverConfig: config.MCPServer{
				Transport: config.MCPTransport{
					Type: "http",
					Port: 8080,
					Path: "/mcp",
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport, err := manager.createHTTPTransport(tt.serverConfig)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, transport)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, transport)
			}
		})
	}
}

func TestClientGetters(t *testing.T) {
	serverConfig := config.MCPServer{
		Name: "test_server",
		Command: []string{"echo", "hello"},
		Transport: config.MCPTransport{
			Type: "stdio",
		},
	}

	client := &Client{
		name:   "test_server",
		config: serverConfig,
	}

	assert.Equal(t, "test_server", client.GetName())
	assert.Equal(t, serverConfig, client.GetConfig())
}
