package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

// OllamaProvider implements the Provider interface for Ollama
type OllamaProvider struct {
	endpoint string
	model    string
	client   *http.Client
	logger   *logrus.Entry
}

// OllamaRequest represents an Ollama API request
type OllamaRequest struct {
	Model    string                   `json:"model"`
	Messages []OllamaMessage          `json:"messages"`
	Stream   bool                     `json:"stream"`
	Tools    []map[string]interface{} `json:"tools,omitempty"`
	Options  map[string]interface{}   `json:"options,omitempty"`
}

// OllamaMessage represents an Ollama message
type OllamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OllamaResponse represents an Ollama API response
type OllamaResponse struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done bool `json:"done"`
}

// NewOllamaProvider creates a new Ollama provider
func NewOllamaProvider(endpoint, model string) *OllamaProvider {
	if endpoint == "" {
		endpoint = "http://localhost:11434"
	}
	if model == "" {
		model = "llama2"
	}

	return &OllamaProvider{
		endpoint: endpoint,
		model:    model,
		client:   &http.Client{},
		logger:   logrus.WithField("provider", "ollama"),
	}
}

// Name returns the provider name
func (p *OllamaProvider) Name() string {
	return "ollama"
}

// ValidateConfig validates the Ollama configuration
func (p *OllamaProvider) ValidateConfig() error {
	// Test connection to Ollama
	resp, err := p.client.Get(p.endpoint + "/api/tags")
	if err != nil {
		return fmt.Errorf("failed to connect to Ollama at %s: %w", p.endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ollama server returned status %d", resp.StatusCode)
	}

	return nil
}

// Chat sends a chat request to Ollama
func (p *OllamaProvider) Chat(ctx context.Context, request *ChatRequest) (*ChatResponse, error) {
	ollamaReq := p.convertRequest(request)
	ollamaReq.Stream = false

	reqBody, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.endpoint+"/api/chat", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &ChatResponse{
		Content:  ollamaResp.Message.Content,
		Finished: ollamaResp.Done,
	}, nil
}

// ChatStream sends a streaming chat request to Ollama
func (p *OllamaProvider) ChatStream(ctx context.Context, request *ChatRequest) (<-chan *ChatResponse, error) {
	ollamaReq := p.convertRequest(request)
	ollamaReq.Stream = true

	reqBody, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.endpoint+"/api/chat", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("Ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

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
				p.logger.WithError(err).Error("Failed to decode streaming response")
				responseChan <- &ChatResponse{
					Error: fmt.Sprintf("decode error: %v", err),
				}
				return
			}

			response := &ChatResponse{
				Content:  ollamaResp.Message.Content,
				Finished: ollamaResp.Done,
			}

			select {
			case responseChan <- response:
			case <-ctx.Done():
				return
			}

			if ollamaResp.Done {
				break
			}
		}
	}()

	return responseChan, nil
}

// convertRequest converts a generic ChatRequest to an Ollama-specific request
func (p *OllamaProvider) convertRequest(request *ChatRequest) *OllamaRequest {
	ollamaReq := &OllamaRequest{
		Model:    p.model,
		Messages: make([]OllamaMessage, 0, len(request.Messages)),
		Options:  make(map[string]interface{}),
	}

	// Add system prompt if provided
	if request.SystemPrompt != "" {
		ollamaReq.Messages = append(ollamaReq.Messages, OllamaMessage{
			Role:    "system",
			Content: request.SystemPrompt,
		})
	}

	// Convert messages
	for _, msg := range request.Messages {
		ollamaReq.Messages = append(ollamaReq.Messages, OllamaMessage(msg))
	}

	// Set temperature if provided
	if request.Temperature != nil {
		ollamaReq.Options["temperature"] = *request.Temperature
	}

	// Convert tools to Ollama format (if supported)
	if len(request.Tools) > 0 {
		tools := make([]map[string]interface{}, len(request.Tools))
		for i, tool := range request.Tools {
			tools[i] = map[string]interface{}{
				"type": "function",
				"function": map[string]interface{}{
					"name":        tool.Name,
					"description": tool.Description,
					"parameters":  tool.Parameters,
				},
			}
		}
		ollamaReq.Tools = tools
	}

	return ollamaReq
}

// Close closes any resources used by the provider
func (p *OllamaProvider) Close() error {
	// HTTP client doesn't need explicit closing
	return nil
}
