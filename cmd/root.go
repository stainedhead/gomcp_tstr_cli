package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	serverName   string
	providerName string
	logLevel     string
	useAllMCP    bool
	logToFile    bool
	jsonRaw      bool
	version      = "1.0.0"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mcp_tstr",
	Short: "MCP Testing CLI - A tool for testing and interacting with MCP servers",
	Long: `mcp_tstr is a comprehensive CLI tool for testing Model Context Protocol (MCP) servers.
It provides capabilities to discover server features, execute tools, and chat with AI models
that have access to MCP server tools and resources.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initLogging()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./mcp_tstr.config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&serverName, "server", "s", "", "MCP server name to interact with")
	rootCmd.PersistentFlags().StringVarP(&providerName, "provider-name", "p", "", "model provider and protocol to use in chat")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "logging level (debug, info, warn, error)")
	rootCmd.PersistentFlags().BoolVarP(&useAllMCP, "use-all-mcp", "u", false, "include all servers in chat session")
	rootCmd.PersistentFlags().BoolVarP(&logToFile, "log-to-file", "f", false, "store logs to persistent file")
	rootCmd.PersistentFlags().BoolVarP(&jsonRaw, "json-raw", "j", false, "turn off json formatting in discovery results")

	// Version flag
	rootCmd.Flags().BoolP("version", "v", false, "show version information")
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		if versionFlag, _ := cmd.Flags().GetBool("version"); versionFlag {
			fmt.Printf("mcp_tstr version %s\n", version)
			return
		}
		_ = cmd.Help()
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("mcp_tstr.config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		logrus.WithField("config", viper.ConfigFileUsed()).Debug("Using config file")
	}
}

// initLogging configures the logging system
func initLogging() {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.WithError(err).Warn("Invalid log level, using info")
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	if logToFile {
		file, err := os.OpenFile("mcp_tstr.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logrus.WithError(err).Warn("Failed to open log file, using stdout")
		} else {
			logrus.SetOutput(file)
		}
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
