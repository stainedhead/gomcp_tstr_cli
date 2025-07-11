package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// showMcpJsonCmd represents the show-mcp-json command
var showMcpJsonCmd = &cobra.Command{
	Use:   "show-mcp-json",
	Short: "Display the complete mcp.json file contents",
	Long: `Display the complete contents of the mcp.json configuration file.
Results are formatted as prettified JSON unless the --json-raw flag is used.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runShowMcpJson()
	},
}

func init() {
	rootCmd.AddCommand(showMcpJsonCmd)
}

func runShowMcpJson() error {
	// Read the mcp.json file directly
	data, err := os.ReadFile("mcp.json")
	if err != nil {
		return fmt.Errorf("failed to read mcp.json: %w", err)
	}

	if jsonRaw {
		// Output raw JSON
		fmt.Println(string(data))
		return nil
	}

	// Parse and re-format for pretty printing
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("failed to parse mcp.json: %w", err)
	}

	prettyData, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format JSON: %w", err)
	}

	fmt.Println(string(prettyData))
	return nil
}
