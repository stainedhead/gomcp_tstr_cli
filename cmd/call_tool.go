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

var (
	toolName   string
	toolParams string
)

// callToolCmd represents the call-tool command
var callToolCmd = &cobra.Command{
	Use:   "call-tool",
	Short: "Execute a specific tool with parameters",
	Long: `Execute a tool provided by the MCP server with the specified parameters.
Parameters should be provided as JSON-RPC formatted string.

Example:
  mcp_tstr call-tool --name "get_weather" --params '{"location":"New York"}'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCallTool()
	},
}

func init() {
	rootCmd.AddCommand(callToolCmd)
	
	callToolCmd.Flags().StringVarP(&toolName, "name", "n", "", "tool name to execute (required)")
	callToolCmd.Flags().StringVarP(&toolParams, "params", "p", "{}", "JSON-RPC formatted parameters")
	_ = callToolCmd.MarkFlagRequired("name")
}

func runCallTool() error {
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

	// Parse tool parameters
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(toolParams), &params); err != nil {
		return fmt.Errorf("failed to parse tool parameters: %w", err)
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

	// Execute the tool
	logrus.WithFields(logrus.Fields{
		"tool":   toolName,
		"params": params,
		"server": targetServer,
	}).Info("Executing tool")

	result, err := client.CallTool(ctx, toolName, params)
	if err != nil {
		return fmt.Errorf("failed to call tool %s: %w", toolName, err)
	}

	// Output results
	return outputJSON(result)
}
