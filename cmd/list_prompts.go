package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sirupsen/logrus"

	"mcp_tstr/internal/config"
	"mcp_tstr/internal/mcp"
)

// listPromptsCmd represents the list-prompts command
var listPromptsCmd = &cobra.Command{
	Use:   "list-prompts",
	Short: "List prompts provided by MCP server",
	Long: `List all prompts provided by the specified MCP server.
Results are formatted as prettified JSON unless the --json-raw flag is used.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runListPrompts()
	},
}

func init() {
	rootCmd.AddCommand(listPromptsCmd)
}

func runListPrompts() error {
	// Load configurations
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	mcpConfig, err := config.LoadMCPConfig()
	if err != nil {
		return fmt.Errorf("failed to load MCP config: %w", err)
	}

	// Determine which server to use
	targetServer := serverName
	if targetServer == "" {
		targetServer = cfg.DefaultServer
	}
	if targetServer == "" {
		return fmt.Errorf("no server specified and no default server configured")
	}

	// Initialize MCP manager
	manager := mcp.NewManager(logrus.StandardLogger())
	defer manager.Close()

	// Initialize the target server
	if err := manager.InitializeServers(mcpConfig, []string{targetServer}); err != nil {
		return fmt.Errorf("failed to initialize MCP servers: %w", err)
	}

	client, err := manager.GetClient(targetServer)
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx := context.Background()

	// Get prompts
	prompts, err := client.ListPrompts(ctx)
	if err != nil {
		return fmt.Errorf("failed to list prompts: %w", err)
	}

	// Output results
	return outputJSON(prompts.Prompts)
}
