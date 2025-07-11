package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/sirupsen/logrus"

	"mcp_tstr/internal/config"
	"mcp_tstr/internal/mcp"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Send a ping request to the MCP server",
	Long: `Send a ping request to the specified MCP server to test connectivity.
Returns the ping result and connection status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPing()
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}

func runPing() error {
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

	// Measure ping time
	start := time.Now()
	err = client.Ping(ctx)
	duration := time.Since(start)

	result := map[string]interface{}{
		"server":   targetServer,
		"success":  err == nil,
		"duration": duration.String(),
	}

	if err != nil {
		result["error"] = err.Error()
		logrus.WithError(err).Errorf("Ping to server %s failed", targetServer)
	} else {
		logrus.Infof("Ping to server %s successful (%s)", targetServer, duration)
	}

	return outputJSON(result)
}
