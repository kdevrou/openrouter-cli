package cli

import (
	"os"

	"github.com/kdevrou/openrouter-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	configPath string
	apiKey     string
	debug      bool
)

// RootCmd is the root command
var RootCmd = &cobra.Command{
	Use:   "openrouter",
	Short: "OpenRouter CLI - Access 400+ AI models from your terminal",
	Long: `OpenRouter CLI is a command-line interface for OpenRouter.ai

It allows you to:
- Send chat completions to 400+ AI models
- List available models with pricing and capabilities
- Pipe text input and output for integration with other tools

Get started:
  openrouter chat "Hello, world!"
  openrouter list
  echo "Tell me a joke" | openrouter chat`,
	Version: "0.1.0",
	Run: func(cmd *cobra.Command, args []string) {
		// Show help if no subcommand
		cmd.Help()
	},
}

// Execute executes the root command
func Execute() error {
	return RootCmd.Execute()
}

func init() {
	// Global flags
	RootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to config file")
	RootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "OpenRouter API key (overrides config)")
	RootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug output")

	// Register subcommands
	RootCmd.AddCommand(chatCmd)
	RootCmd.AddCommand(listCmd)
	RootCmd.AddCommand(configCmd)
}

// GetConfig loads the configuration with command-line overrides
func GetConfig() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Override with command-line flags
	if apiKey != "" {
		cfg.APIKey = apiKey
	}

	// Validate API key is set
	if cfg.APIKey == "" {
		return nil, config.ErrNoAPIKey
	}

	return cfg, nil
}

// PrintSetupError prints an error message when API key is not configured
func PrintSetupError() {
	PrintError("No API key found")
	PrintSetupInstructions()
	os.Exit(1)
}
