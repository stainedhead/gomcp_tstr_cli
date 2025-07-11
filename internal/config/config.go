package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/viper"

	"mcp_tstr/internal/constants"
)

// Config represents the application configuration
type Config struct {
	DefaultServer   string                 `yaml:"default_server" mapstructure:"default_server"`
	DefaultProvider string                 `yaml:"default_provider" mapstructure:"default_provider"`
	DefaultModel    string                 `yaml:"default_model" mapstructure:"default_model"`
	Providers       map[string]interface{} `yaml:"providers" mapstructure:"providers"`
	Logging         LoggingConfig          `yaml:"logging" mapstructure:"logging"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level  string `yaml:"level" mapstructure:"level"`
	ToFile bool   `yaml:"to_file" mapstructure:"to_file"`
}

// ProviderConfig represents a generic provider configuration
type ProviderConfig struct {
	Type     string                 `yaml:"type" mapstructure:"type"`
	Endpoint string                 `yaml:"endpoint" mapstructure:"endpoint"`
	APIKey   string                 `yaml:"api_key" mapstructure:"api_key"`
	Model    string                 `yaml:"model" mapstructure:"model"`
	Extra    map[string]interface{} `yaml:"extra" mapstructure:"extra"`
}

// OllamaConfig represents Ollama provider configuration
type OllamaConfig struct {
	Endpoint string `yaml:"endpoint" mapstructure:"endpoint"`
	Model    string `yaml:"model" mapstructure:"model"`
}

// OpenAIConfig represents OpenAI provider configuration
type OpenAIConfig struct {
	APIKey      string `yaml:"api_key" mapstructure:"api_key"`
	Model       string `yaml:"model" mapstructure:"model"`
	BaseURL     string `yaml:"base_url" mapstructure:"base_url"`
	Temperature string `yaml:"temperature" mapstructure:"temperature"`
}

// AWSBedrockConfig represents AWS Bedrock provider configuration
type AWSBedrockConfig struct {
	Region          string `yaml:"region" mapstructure:"region"`
	AccessKeyID     string `yaml:"access_key_id" mapstructure:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key" mapstructure:"secret_access_key"`
	Model           string `yaml:"model" mapstructure:"model"`
}

// GoogleAIConfig represents Google AI provider configuration
type GoogleAIConfig struct {
	APIKey string `yaml:"api_key" mapstructure:"api_key"`
	Model  string `yaml:"model" mapstructure:"model"`
}

// AnthropicConfig represents Anthropic provider configuration
type AnthropicConfig struct {
	APIKey string `yaml:"api_key" mapstructure:"api_key"`
	Model  string `yaml:"model" mapstructure:"model"`
}

// MCPServer represents an MCP server configuration
type MCPServer struct {
	Name      string                 `json:"name"`
	Command   []string               `json:"command,omitempty"`
	Args      []string               `json:"args,omitempty"`
	Env       map[string]string      `json:"env,omitempty"`
	Transport MCPTransport           `json:"transport"`
	Extra     map[string]interface{} `json:"extra,omitempty"`
}

// MCPTransport represents the transport configuration for an MCP server
type MCPTransport struct {
	Type string `json:"type"` // "stdio", "http", "sse"
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
	Path string `json:"path,omitempty"`
}

// MCPConfig represents the MCP servers configuration
type MCPConfig struct {
	Servers map[string]MCPServer `json:"servers"`
}

// Load loads the application configuration
func Load() (*Config, error) {
	var cfg Config

	// Set defaults
	viper.SetDefault("default_server", "")
	viper.SetDefault("default_provider", constants.DefaultProvider)
	viper.SetDefault("default_model", constants.DefaultModel)
	viper.SetDefault("logging.level", constants.DefaultLogLevel)
	viper.SetDefault("logging.to_file", false)

	// Set up environment variable mappings
	_ = viper.BindEnv("providers.ollama.endpoint", "OLLAMA_ENDPOINT")
	_ = viper.BindEnv("providers.ollama.model", "OLLAMA_MODEL")
	_ = viper.BindEnv("providers.openai.api_key", "OPENAI_API_KEY")
	_ = viper.BindEnv("providers.openai.model", "OPENAI_MODEL")
	_ = viper.BindEnv("providers.aws_bedrock.region", "AWS_REGION")
	_ = viper.BindEnv("providers.aws_bedrock.access_key_id", "AWS_ACCESS_KEY_ID")
	_ = viper.BindEnv("providers.aws_bedrock.secret_access_key", "AWS_SECRET_ACCESS_KEY")
	_ = viper.BindEnv("providers.google_ai.api_key", "GOOGLE_AI_API_KEY")
	_ = viper.BindEnv("providers.anthropic.api_key", "ANTHROPIC_API_KEY")

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// LoadMCPConfig loads the MCP servers configuration from mcp.json
func LoadMCPConfig() (*MCPConfig, error) {
	file, err := os.ReadFile(constants.MCPConfigFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", constants.MCPConfigFileName, err)
	}

	var mcpConfig MCPConfig
	if err := json.Unmarshal(file, &mcpConfig); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", constants.MCPConfigFileName, err)
	}

	return &mcpConfig, nil
}

// GetProviderConfig extracts provider-specific configuration
func (c *Config) GetProviderConfig(providerName string) (interface{}, error) {
	if c.Providers == nil {
		return nil, fmt.Errorf("no providers configured")
	}

	providerData, exists := c.Providers[providerName]
	if !exists {
		return nil, fmt.Errorf("provider %s not configured", providerName)
	}

	return providerData, nil
}
