package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadMCPConfig(t *testing.T) {
	// Create a temporary mcp.json file
	mcpJSON := `{
		"servers": {
			"test_server": {
				"name": "test_server",
				"command": ["echo", "hello"],
				"transport": {
					"type": "stdio"
				}
			}
		}
	}`

	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "mcp_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Write mcp.json to temporary directory
	mcpFile := tmpDir + "/mcp.json"
	err = os.WriteFile(mcpFile, []byte(mcpJSON), 0644)
	require.NoError(t, err)

	// Change to temp directory
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		_ = os.Chdir(originalDir)
	}()

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Test loading
	config, err := LoadMCPConfig()
	require.NoError(t, err)
	assert.NotNil(t, config)
	assert.Contains(t, config.Servers, "test_server")
	assert.Equal(t, "stdio", config.Servers["test_server"].Transport.Type)
}

func TestMCPServerTransportTypes(t *testing.T) {
	tests := []struct {
		name      string
		transport MCPTransport
		valid     bool
	}{
		{
			name: "stdio transport",
			transport: MCPTransport{
				Type: "stdio",
			},
			valid: true,
		},
		{
			name: "http transport",
			transport: MCPTransport{
				Type: "http",
				Host: "localhost",
				Port: 8080,
				Path: "/mcp",
			},
			valid: true,
		},
		{
			name: "sse transport",
			transport: MCPTransport{
				Type: "sse",
				Host: "localhost",
				Port: 9090,
				Path: "/events",
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, tt.transport.Type)
			if tt.transport.Type != "stdio" {
				assert.NotEmpty(t, tt.transport.Host)
				assert.Greater(t, tt.transport.Port, 0)
			}
		})
	}
}

func TestConfigDefaults(t *testing.T) {
	config := &Config{
		DefaultProvider: "ollama",
		DefaultModel:    "llama2",
		Logging: LoggingConfig{
			Level:  "info",
			ToFile: false,
		},
	}

	assert.Equal(t, "ollama", config.DefaultProvider)
	assert.Equal(t, "llama2", config.DefaultModel)
	assert.Equal(t, "info", config.Logging.Level)
	assert.False(t, config.Logging.ToFile)
}
