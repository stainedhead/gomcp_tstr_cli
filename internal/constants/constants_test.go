package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstants(t *testing.T) {
	// Test that all constants are properly defined
	assert.Equal(t, "mcp_tstr", AppName)
	assert.Equal(t, "1.0.0", AppVersion)
	assert.Equal(t, "mcp_tstr.config", ConfigFileName)
	assert.Equal(t, "mcp.json", MCPConfigFileName)
	assert.Equal(t, "info", DefaultLogLevel)
	assert.Equal(t, "ollama", DefaultProvider)
	assert.Equal(t, "llama2", DefaultModel)
}

func TestConstantsNotEmpty(t *testing.T) {
	// Ensure no constants are empty strings
	assert.NotEmpty(t, AppName)
	assert.NotEmpty(t, AppVersion)
	assert.NotEmpty(t, ConfigFileName)
	assert.NotEmpty(t, MCPConfigFileName)
	assert.NotEmpty(t, DefaultLogLevel)
	assert.NotEmpty(t, DefaultProvider)
	assert.NotEmpty(t, DefaultModel)
}

func TestConfigFileExtension(t *testing.T) {
	// Ensure config file doesn't have .yaml extension
	assert.NotContains(t, ConfigFileName, ".yaml")
	assert.NotContains(t, ConfigFileName, ".yml")
	
	// Ensure MCP config file has .json extension
	assert.Contains(t, MCPConfigFileName, ".json")
}
