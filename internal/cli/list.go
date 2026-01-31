package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/kdevrou/openrouter-cli/internal/api"
)

var (
	// List command flags
	filterName string
	jsonList   bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available models with pricing and capabilities",
	Long: `Display all available models on OpenRouter with pricing and capabilities.

Use --filter to search for specific models:
  openrouter list --filter gpt
  openrouter list --filter claude

Use --json to get raw JSON output for scripting:
  openrouter list --json | jq '.[] | .id'`,

	RunE: runList,
}

func runList(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := GetConfig()
	if err != nil {
		PrintSetupError()
	}

	// Create API client
	apiClient := api.NewClient(cfg.APIBaseURL, cfg.APIKey, cfg.Timeout)

	// Fetch models
	if debug {
		fmt.Fprintf(os.Stderr, "Fetching models from %s\n", cfg.APIBaseURL)
	}

	models, err := apiClient.ListModels()
	if err != nil {
		if apiErr, ok := err.(*api.APIError); ok {
			PrintAPIError(apiErr)
		} else {
			PrintError(err.Error())
		}
		return err
	}

	// Filter models if requested
	if filterName != "" {
		filtered := make([]api.Model, 0)
		filterLower := strings.ToLower(filterName)
		for _, m := range models {
			if strings.Contains(strings.ToLower(m.ID), filterLower) ||
				strings.Contains(strings.ToLower(m.Name), filterLower) {
				filtered = append(filtered, m)
			}
		}
		models = filtered
	}

	if len(models) == 0 {
		PrintError("No models found")
		if filterName != "" {
			fmt.Fprintf(os.Stderr, "Try searching without filters or with different keywords\n")
		}
		return nil
	}

	// Format output
	format := FormatPretty
	if jsonList {
		format = FormatJSON
	}

	return FormatModelList(models, format)
}

func init() {
	listCmd.Flags().StringVar(&filterName, "filter", "", "Filter models by name or ID")
	listCmd.Flags().BoolVar(&jsonList, "json", false, "Output as JSON")
}
