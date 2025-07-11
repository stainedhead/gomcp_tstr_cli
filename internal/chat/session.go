package chat

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"mcp_tstr/internal/mcp"
	"mcp_tstr/internal/providers"
)

// Session represents a chat session
type Session struct {
	provider    providers.Provider
	mcpManager  *mcp.Manager
	messages    []providers.Message
	tools       []providers.Tool
	systemPrompt string
	logger      *logrus.Entry
}

// NewSession creates a new chat session
func NewSession(provider providers.Provider, mcpManager *mcp.Manager) *Session {
	return &Session{
		provider:   provider,
		mcpManager: mcpManager,
		messages:   make([]providers.Message, 0),
		tools:      make([]providers.Tool, 0),
		systemPrompt: `You are a helpful AI assistant with access to various tools through MCP (Model Context Protocol) servers. 
You can use these tools to help users with their requests. When you need to use a tool, make sure to call it with the appropriate parameters.
Be helpful, accurate, and explain what you're doing when using tools.`,
		logger: logrus.WithField("component", "chat"),
	}
}

// SetSystemPrompt sets the system prompt for the chat session
func (s *Session) SetSystemPrompt(prompt string) {
	s.systemPrompt = prompt
}

// LoadTools loads available tools from MCP servers
func (s *Session) LoadTools(ctx context.Context) error {
	s.tools = make([]providers.Tool, 0)

	for _, client := range s.mcpManager.GetAllClients() {
		toolsResult, err := client.ListTools(ctx)
		if err != nil {
			s.logger.WithError(err).Warnf("Failed to load tools from server %s", client.GetName())
			continue
		}

		for _, tool := range toolsResult.Tools {
			// Convert schema to map if available
			var parameters map[string]interface{}
			if tool.InputSchema != nil {
				// For now, we'll create a simple representation
				// In a full implementation, you'd properly convert the JSON schema
				parameters = map[string]interface{}{
					"type": "object",
					"description": "Tool parameters",
				}
			}

			s.tools = append(s.tools, providers.Tool{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  parameters,
			})
		}

		s.logger.Infof("Loaded %d tools from server %s", len(toolsResult.Tools), client.GetName())
	}

	s.logger.Infof("Total tools available: %d", len(s.tools))
	return nil
}

// Start starts an interactive chat session
func (s *Session) Start(ctx context.Context) error {
	fmt.Println("Starting chat session. Type 'bye', 'exit', 'end', or 'quit' to end the session.")
	fmt.Println("Available tools:", len(s.tools))
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("You: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// Check for exit commands
		if isExitCommand(input) {
			fmt.Println("Goodbye!")
			break
		}

		// Add user message
		s.messages = append(s.messages, providers.Message{
			Role:    "user",
			Content: input,
		})

		// Send chat request
		if err := s.processMessage(ctx); err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
	}

	return scanner.Err()
}

// processMessage processes a user message and generates a response
func (s *Session) processMessage(ctx context.Context) error {
	request := &providers.ChatRequest{
		Messages:     s.messages,
		Tools:        s.tools,
		Stream:       true,
		SystemPrompt: s.systemPrompt,
	}

	// Use streaming for better user experience
	responseChan, err := s.provider.ChatStream(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to start chat stream: %w", err)
	}

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
					s.logger.WithError(err).Errorf("Failed to handle tool call: %s", toolCall.Name)
					fmt.Printf("\n[Tool call failed: %v]", err)
				}
			}
		}

		if response.Finished {
			break
		}
	}

	fmt.Println() // New line after response

	// Add assistant response to conversation history
	if fullResponse.Len() > 0 {
		s.messages = append(s.messages, providers.Message{
			Role:    "assistant",
			Content: fullResponse.String(),
		})
	}

	return nil
}

// handleToolCall executes a tool call through the appropriate MCP server
func (s *Session) handleToolCall(ctx context.Context, toolCall providers.ToolCall) error {
	s.logger.WithFields(logrus.Fields{
		"tool": toolCall.Name,
		"args": toolCall.Arguments,
	}).Info("Executing tool call")

	// Find which server has this tool
	var targetClient *mcp.Client
	for _, client := range s.mcpManager.GetAllClients() {
		toolsResult, err := client.ListTools(ctx)
		if err != nil {
			continue
		}

		for _, tool := range toolsResult.Tools {
			if tool.Name == toolCall.Name {
				targetClient = client
				break
			}
		}

		if targetClient != nil {
			break
		}
	}

	if targetClient == nil {
		return fmt.Errorf("tool %s not found in any connected server", toolCall.Name)
	}

	// Execute the tool
	result, err := targetClient.CallTool(ctx, toolCall.Name, toolCall.Arguments)
	if err != nil {
		return fmt.Errorf("failed to call tool %s: %w", toolCall.Name, err)
	}

	// Display tool result
	fmt.Printf("\n[Tool %s executed successfully]", toolCall.Name)
	if len(result.Content) > 0 {
		fmt.Printf("\n[Result: %v]", result.Content)
	}

	return nil
}

// isExitCommand checks if the input is an exit command
func isExitCommand(input string) bool {
	lower := strings.ToLower(input)
	return lower == "bye" || lower == "exit" || lower == "end" || lower == "quit"
}
