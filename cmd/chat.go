package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sirupsen/logrus"

	"mcp_tstr/internal/chat"
	"mcp_tstr/internal/config"
	"mcp_tstr/internal/mcp"
	"mcp_tstr/internal/providers"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start an interactive chat session with AI model and MCP tools",
	Long: `Start an interactive chat session where you can chat with an AI model that has access
to tools provided by MCP servers. The model can use these tools to help answer your questions
and perform tasks.

The chat session will continue until you type 'bye', 'exit', 'end', or 'quit'.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runChat()
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}

func runChat() error {
	// Load configurations
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	mcpConfig, err := config.LoadMCPConfig()
	if err != nil {
		return fmt.Errorf("failed to load MCP config: %w", err)
	}

	// Determine which provider to use
	targetProvider := providerName
	if targetProvider == "" {
		targetProvider = cfg.DefaultProvider
	}
	if targetProvider == "" {
		return fmt.Errorf("no provider specified and no default provider configured")
	}

	// Create provider
	provider, err := providers.NewProvider(targetProvider, cfg)
	if err != nil {
		return fmt.Errorf("failed to create provider %s: %w", targetProvider, err)
	}
	defer provider.Close()

	// Determine which servers to use
	var serverNames []string
	if useAllMCP {
		// Use all configured servers
		for name := range mcpConfig.Servers {
			serverNames = append(serverNames, name)
		}
	} else {
		// Use specified server or default
		targetServer := serverName
		if targetServer == "" {
			targetServer = cfg.DefaultServer
		}
		if targetServer == "" {
			return fmt.Errorf("no server specified and no default server configured")
		}
		serverNames = []string{targetServer}
	}

	// Initialize MCP manager
	manager := mcp.NewManager(logrus.StandardLogger())
	defer manager.Close()

	// Initialize MCP servers
	if err := manager.InitializeServers(mcpConfig, serverNames); err != nil {
		return fmt.Errorf("failed to initialize MCP servers: %w", err)
	}

	// Create chat session
	session := chat.NewSession(provider, manager)

	// Load available tools from MCP servers
	ctx := context.Background()
	if err := session.LoadTools(ctx); err != nil {
		logrus.WithError(err).Warn("Failed to load some tools, continuing anyway")
	}

	// Start interactive chat
	fmt.Printf("Chat session started with provider: %s\n", provider.Name())
	fmt.Printf("Connected MCP servers: %d\n", len(manager.GetAllClients()))
	fmt.Println()

	return session.Start(ctx)
}
