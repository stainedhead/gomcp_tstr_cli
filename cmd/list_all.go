package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sirupsen/logrus"

	"mcp_tstr/internal/config"
	"mcp_tstr/internal/mcp"
)

// listAllCmd represents the list-all command
var listAllCmd = &cobra.Command{
	Use:   "list-all",
	Short: "List all tools, resources, and prompts from MCP server",
	Long: `List all capabilities (tools, resources, and prompts) provided by the specified MCP server.
Results are formatted as prettified JSON unless the --json-raw flag is used.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runListAll()
	},
}

func init() {
	rootCmd.AddCommand(listAllCmd)
}

func runListAll() error {
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

	// Collect all capabilities
	result := make(map[string]interface{})

	// Get tools
	if tools, err := client.ListTools(ctx); err != nil {
		logrus.WithError(err).Warn("Failed to list tools")
		result["tools"] = map[string]string{"error": err.Error()}
	} else {
		result["tools"] = tools.Tools
	}

	// Get resources
	if resources, err := client.ListResources(ctx); err != nil {
		logrus.WithError(err).Warn("Failed to list resources")
		result["resources"] = map[string]string{"error": err.Error()}
	} else {
		result["resources"] = resources.Resources
	}

	// Get prompts
	if prompts, err := client.ListPrompts(ctx); err != nil {
		logrus.WithError(err).Warn("Failed to list prompts")
		result["prompts"] = map[string]string{"error": err.Error()}
	} else {
		result["prompts"] = prompts.Prompts
	}

	// Format and output results
	return outputJSON(result)
}

func outputJSON(data interface{}) error {
	var output []byte
	var err error

	if jsonRaw {
		output, err = json.Marshal(data)
	} else {
		output, err = json.MarshalIndent(data, "", "  ")
	}

	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Println(string(output))
	return nil
}
