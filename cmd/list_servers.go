package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"mcp_tstr/internal/config"
)

// listServersCmd represents the list-servers command
var listServersCmd = &cobra.Command{
	Use:   "list-servers",
	Short: "List configured MCP servers",
	Long: `List all MCP servers configured in mcp.json with their protocol and command information.
Results are formatted as prettified JSON unless the --json-raw flag is used.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runListServers()
	},
}

func init() {
	rootCmd.AddCommand(listServersCmd)
}

func runListServers() error {
	mcpConfig, err := config.LoadMCPConfig()
	if err != nil {
		return fmt.Errorf("failed to load MCP config: %w", err)
	}

	// Create a simplified view of servers for output
	servers := make(map[string]interface{})
	for name, server := range mcpConfig.Servers {
		serverInfo := map[string]interface{}{
			"name":      name,
			"transport": server.Transport,
		}
		
		if len(server.Command) > 0 {
			serverInfo["command"] = server.Command
		}
		
		if len(server.Args) > 0 {
			serverInfo["args"] = server.Args
		}
		
		if len(server.Env) > 0 {
			serverInfo["env"] = server.Env
		}
		
		servers[name] = serverInfo
	}

	return outputJSON(servers)
}
