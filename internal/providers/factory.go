package providers

import (
	"fmt"

	"mcp_tstr/internal/config"
)

// NewProvider creates a new provider based on the configuration
func NewProvider(providerName string, cfg *config.Config) (Provider, error) {
	providerConfig, err := cfg.GetProviderConfig(providerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider config: %w", err)
	}

	switch providerName {
	case "ollama":
		return createOllamaProvider(providerConfig)
	case "openai":
		return createOpenAIProvider(providerConfig)
	case "aws_bedrock":
		return createAWSBedrockProvider(providerConfig)
	case "google_ai":
		return createGoogleAIProvider(providerConfig)
	case "anthropic":
		return createAnthropicProvider(providerConfig)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", providerName)
	}
}

func createOllamaProvider(configData interface{}) (Provider, error) {
	configMap, ok := configData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid ollama configuration format")
	}

	endpoint, _ := configMap["endpoint"].(string)
	model, _ := configMap["model"].(string)

	provider := NewOllamaProvider(endpoint, model)
	return provider, provider.ValidateConfig()
}

func createOpenAIProvider(configData interface{}) (Provider, error) {
	// TODO: Implement OpenAI provider
	return nil, fmt.Errorf("OpenAI provider not yet implemented")
}

func createAWSBedrockProvider(configData interface{}) (Provider, error) {
	// TODO: Implement AWS Bedrock provider
	return nil, fmt.Errorf("AWS Bedrock provider not yet implemented")
}

func createGoogleAIProvider(configData interface{}) (Provider, error) {
	// TODO: Implement Google AI provider
	return nil, fmt.Errorf("Google AI provider not yet implemented")
}

func createAnthropicProvider(configData interface{}) (Provider, error) {
	// TODO: Implement Anthropic provider
	return nil, fmt.Errorf("Anthropic provider not yet implemented")
}
