package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/yourusername/openrouter-cli/internal/api"
)

// OutputFormat represents the desired output format
type OutputFormat string

const (
	FormatPretty OutputFormat = "pretty"
	FormatRaw    OutputFormat = "raw"
	FormatJSON   OutputFormat = "json"
)

// FormatChatResponse formats a chat completion response
func FormatChatResponse(resp *api.ChatCompletionResponse, format OutputFormat) error {
	if len(resp.Choices) == 0 {
		return fmt.Errorf("no choices in response")
	}

	choice := resp.Choices[0]
	message := choice.Message.Content

	switch format {
	case FormatRaw:
		fmt.Print(message)
	case FormatJSON:
		data, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal response: %w", err)
		}
		fmt.Println(string(data))
	default: // FormatPretty
		fmt.Printf("%s\n", message)

		// Print usage stats
		if resp.Usage.TotalTokens > 0 {
			fmt.Printf("\n%s\n",
				color.CyanString(fmt.Sprintf("Tokens used: %d (prompt: %d, completion: %d)",
					resp.Usage.TotalTokens,
					resp.Usage.PromptTokens,
					resp.Usage.CompletionTokens)))
		}
	}

	return nil
}

// FormatModelList formats a list of models as a table
func FormatModelList(models []api.Model, format OutputFormat) error {
	if format == FormatJSON {
		data, err := json.MarshalIndent(models, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal models: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Format as table for pretty/raw output
	// Print header
	fmt.Printf("%-35s | %-15s | %-30s | %-10s\n",
		"Model ID", "Context Length", "Pricing (in/out)", "Modality")
	fmt.Println(strings.Repeat("-", 95))

	for _, model := range models {
		// Shorten model ID for readability
		modelID := model.ID
		if len(modelID) > 33 {
			modelID = modelID[:30] + "..."
		}

		// Format pricing (remove leading zeros and format)
		pricing := formatPrice(model.Pricing.Prompt) + " / " + formatPrice(model.Pricing.Completion)
		if len(pricing) > 28 {
			pricing = pricing[:25] + "..."
		}

		// Format modality
		modality := model.Architecture.Modality
		if modality == "" {
			modality = "text"
		}
		if len(modality) > 8 {
			modality = modality[:5] + "..."
		}

		fmt.Printf("%-35s | %-15d | %-30s | %-10s\n",
			modelID,
			model.ContextLength,
			pricing,
			modality)
	}
	return nil
}

// formatPrice formats a price string for display
func formatPrice(price string) string {
	if price == "" || price == "0" {
		return "free"
	}
	if len(price) > 15 {
		return price[:15] + "..."
	}
	return price
}

// PrintError prints an error in a user-friendly way
func PrintError(message string) {
	fmt.Fprintf(os.Stderr, "%s %s\n",
		color.RedString("Error:"),
		message)
}

// PrintAPIError prints an API error with helpful context
func PrintAPIError(err *api.APIError) {
	if err.Type != "" {
		fmt.Fprintf(os.Stderr, "%s %s (%s, HTTP %d)\n",
			color.RedString("API Error:"),
			err.Message,
			err.Type,
			err.StatusCode)
	} else {
		fmt.Fprintf(os.Stderr, "%s %s (HTTP %d)\n",
			color.RedString("API Error:"),
			err.Message,
			err.StatusCode)
	}
}

// PrintSetupInstructions prints API key setup instructions
func PrintSetupInstructions() {
	fmt.Fprintf(os.Stderr, "\n%s\n", color.YellowString("To set up OpenRouter CLI:"))
	fmt.Fprintf(os.Stderr, "\n1. Get an API key from https://openrouter.ai\n")
	fmt.Fprintf(os.Stderr, "2. Set it using one of:\n")
	fmt.Fprintf(os.Stderr, "   - Environment variable: %s\n", color.CyanString("export OPENROUTER_API_KEY=sk-..."))
	fmt.Fprintf(os.Stderr, "   - Config file: %s\n", color.CyanString("~/.config/openrouter/config.yaml"))
	fmt.Fprintf(os.Stderr, "\n")
}

// WordWrap wraps text to a specified width
func WordWrap(text string, width int) string {
	if width <= 0 {
		return text
	}

	lines := strings.Split(text, "\n")
	var result []string

	for _, line := range lines {
		if len(line) <= width {
			result = append(result, line)
			continue
		}

		words := strings.Fields(line)
		var currentLine string

		for _, word := range words {
			if len(currentLine)+len(word)+1 > width {
				if currentLine != "" {
					result = append(result, currentLine)
				}
				currentLine = word
			} else {
				if currentLine != "" {
					currentLine += " "
				}
				currentLine += word
			}
		}

		if currentLine != "" {
			result = append(result, currentLine)
		}
	}

	return strings.Join(result, "\n")
}
