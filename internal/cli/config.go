package cli

import (
	"fmt"
	"strings"

	"github.com/kdevrou/openrouter-cli/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long: `Manage OpenRouter CLI configuration.

Examples:
  openrouter config get api_key
  openrouter config set default_model openai/gpt-4
  openrouter config add-unavailable qwen/model:free
  openrouter config remove-unavailable qwen/model:free
  openrouter config list-unavailable`,
}

var getCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil && err != config.ErrNoAPIKey {
			PrintError(err.Error())
			return err
		}
		if cfg == nil {
			cfg = config.DefaultConfig()
		}

		key := args[0]
		switch key {
		case "api_key":
			if cfg.APIKey == "" {
				fmt.Println("(not set)")
			} else {
				// Show only last 4 characters for security
				visible := "sk-..." + cfg.APIKey[len(cfg.APIKey)-4:]
				fmt.Println(visible)
			}
		case "default_model":
			fmt.Println(cfg.DefaultModel)
		case "default_temperature":
			fmt.Println(cfg.DefaultTemp)
		case "default_max_tokens":
			fmt.Println(cfg.DefaultMaxTokens)
		case "output_format":
			fmt.Println(cfg.OutputFormat)
		case "api_base_url":
			fmt.Println(cfg.APIBaseURL)
		case "timeout":
			fmt.Println(cfg.Timeout)
		case "unavailable_models":
			if len(cfg.UnavailableModels) == 0 {
				fmt.Println("(none)")
			} else {
				for _, m := range cfg.UnavailableModels {
					fmt.Println(m)
				}
			}
		default:
			PrintError(fmt.Sprintf("unknown config key: %s", key))
			return fmt.Errorf("unknown key")
		}
		return nil
	},
}

var setCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil && err != config.ErrNoAPIKey {
			PrintError(err.Error())
			return err
		}
		if cfg == nil {
			cfg = config.DefaultConfig()
		}

		key := args[0]
		value := args[1]

		switch key {
		case "api_key":
			cfg.APIKey = value
		case "default_model":
			cfg.DefaultModel = value
		case "default_temperature":
			var temp float64
			_, err := fmt.Sscanf(value, "%f", &temp)
			if err != nil {
				PrintError("temperature must be a number")
				return err
			}
			cfg.DefaultTemp = temp
		case "default_max_tokens":
			var tokens int
			_, err := fmt.Sscanf(value, "%d", &tokens)
			if err != nil {
				PrintError("max_tokens must be an integer")
				return err
			}
			cfg.DefaultMaxTokens = tokens
		case "output_format":
			if value != "pretty" && value != "raw" && value != "json" {
				PrintError("output_format must be: pretty, raw, or json")
				return fmt.Errorf("invalid format")
			}
			cfg.OutputFormat = value
		case "timeout":
			var timeout int
			_, err := fmt.Sscanf(value, "%d", &timeout)
			if err != nil {
				PrintError("timeout must be an integer (seconds)")
				return err
			}
			cfg.Timeout = timeout
		default:
			PrintError(fmt.Sprintf("unknown config key: %s", key))
			return fmt.Errorf("unknown key")
		}

		if err := config.Save(cfg); err != nil {
			PrintError(fmt.Sprintf("failed to save config: %v", err))
			return err
		}

		fmt.Printf("✓ Set %s = %s\n", key, value)
		return nil
	},
}

var addUnavailableCmd = &cobra.Command{
	Use:   "add-unavailable <model_id>",
	Short: "Mark a model as unavailable (won't appear in list)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil && err != config.ErrNoAPIKey {
			PrintError(err.Error())
			return err
		}
		if cfg == nil {
			cfg = config.DefaultConfig()
		}

		modelID := args[0]
		if err := cfg.AddUnavailableModel(modelID); err != nil {
			PrintError(err.Error())
			return err
		}

		if err := config.Save(cfg); err != nil {
			PrintError(fmt.Sprintf("failed to save config: %v", err))
			return err
		}

		fmt.Printf("✓ Marked %s as unavailable\n", modelID)
		return nil
	},
}

var removeUnavailableCmd = &cobra.Command{
	Use:   "remove-unavailable <model_id>",
	Short: "Remove a model from the unavailable list",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil && err != config.ErrNoAPIKey {
			PrintError(err.Error())
			return err
		}
		if cfg == nil {
			cfg = config.DefaultConfig()
		}

		modelID := args[0]
		if err := cfg.RemoveUnavailableModel(modelID); err != nil {
			PrintError(err.Error())
			return err
		}

		if err := config.Save(cfg); err != nil {
			PrintError(fmt.Sprintf("failed to save config: %v", err))
			return err
		}

		fmt.Printf("✓ Removed %s from unavailable list\n", modelID)
		return nil
	},
}

var listUnavailableCmd = &cobra.Command{
	Use:   "list-unavailable",
	Short: "List all models marked as unavailable",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil && err != config.ErrNoAPIKey {
			PrintError(err.Error())
			return err
		}
		if cfg == nil {
			cfg = config.DefaultConfig()
		}

		if len(cfg.UnavailableModels) == 0 {
			fmt.Println("No unavailable models configured.")
			return nil
		}

		fmt.Println("Unavailable models (filtered from 'openrouter list'):")
		for i, m := range cfg.UnavailableModels {
			fmt.Printf("  %d. %s\n", i+1, m)
		}
		return nil
	},
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show all configuration settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil && err != config.ErrNoAPIKey {
			PrintError(err.Error())
			return err
		}
		if cfg == nil {
			cfg = config.DefaultConfig()
		}

		fmt.Println("Configuration:")
		fmt.Printf("  API Key: %s\n", maskAPIKey(cfg.APIKey))
		fmt.Printf("  Default Model: %s\n", cfg.DefaultModel)
		fmt.Printf("  Default Temperature: %v\n", cfg.DefaultTemp)
		fmt.Printf("  Default Max Tokens: %d\n", cfg.DefaultMaxTokens)
		fmt.Printf("  Output Format: %s\n", cfg.OutputFormat)
		fmt.Printf("  API Base URL: %s\n", cfg.APIBaseURL)
		fmt.Printf("  Timeout: %d seconds\n", cfg.Timeout)
		fmt.Printf("  Unavailable Models: %d\n", len(cfg.UnavailableModels))
		if len(cfg.UnavailableModels) > 0 {
			fmt.Println("    " + strings.Join(cfg.UnavailableModels, ", "))
		}
		fmt.Printf("  Config File: %s\n", config.GetConfigPath())
		return nil
	},
}

func maskAPIKey(key string) string {
	if key == "" {
		return "(not set)"
	}
	if len(key) <= 4 {
		return "****"
	}
	return "sk-..." + key[len(key)-4:]
}

func init() {
	configCmd.AddCommand(getCmd)
	configCmd.AddCommand(setCmd)
	configCmd.AddCommand(addUnavailableCmd)
	configCmd.AddCommand(removeUnavailableCmd)
	configCmd.AddCommand(listUnavailableCmd)
	configCmd.AddCommand(showCmd)
}
