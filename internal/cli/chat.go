package cli

import (
	"fmt"
	"os"

	"github.com/kdevrou/openrouter-cli/internal/api"
	"github.com/kdevrou/openrouter-cli/internal/util"
	"github.com/spf13/cobra"
)

var (
	// Chat command flags
	model       string
	temperature float64
	maxTokens   int
	rawOutput   bool
	jsonOutput  bool
	useStdin    bool
)

var chatCmd = &cobra.Command{
	Use:   "chat [prompt]",
	Short: "Send a chat completion request",
	Long: `Send a prompt to an AI model via OpenRouter.

Input options:
  openrouter chat "What is Go?"                    # Argument only
  echo "Explain quantum computing" | openrouter chat  # Pipe only
  cat file.txt | openrouter chat --stdin "Analyze:"   # Combine both

Flags let you customize the request:
  -m, --model: Choose which model to use
  -t, --temperature: Adjust response creativity (0.0-2.0)
  --max-tokens: Limit response length
  --stdin: Combine argument with piped input
  --raw: Output only the response text (for piping)
  --json: Output full API response as JSON`,

	Args: cobra.MaximumNArgs(1),
	RunE: runChat,
}

func runChat(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := GetConfig()
	if err != nil {
		PrintSetupError()
	}

	// Get input from args or stdin
	prompt, err := util.CombineInputWithStdin(args, useStdin)
	if err != nil {
		PrintError(err.Error())
		return fmt.Errorf("no input provided")
	}

	if prompt == "" {
		PrintError("prompt cannot be empty")
		return fmt.Errorf("empty prompt")
	}

	// Use provided model or default
	selectedModel := model
	if selectedModel == "" {
		selectedModel = cfg.DefaultModel
	}

	// Use provided temperature or default
	selectedTemp := temperature
	if selectedTemp == 0 && !cmd.Flags().Changed("temperature") {
		selectedTemp = cfg.DefaultTemp
	}

	// Use provided maxTokens or default
	selectedMaxTokens := maxTokens
	if selectedMaxTokens == 0 && !cmd.Flags().Changed("max-tokens") {
		selectedMaxTokens = cfg.DefaultMaxTokens
	}

	// Create API client
	apiClient := api.NewClient(cfg.APIBaseURL, cfg.APIKey, cfg.Timeout)

	// Build request
	chatReq := &api.ChatCompletionRequest{
		Model: selectedModel,
		Messages: []api.Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: selectedTemp,
		MaxTokens:   selectedMaxTokens,
	}

	// Send request
	if debug {
		fmt.Fprintf(os.Stderr, "Sending request to %s with model: %s\n", cfg.APIBaseURL, selectedModel)
	}

	resp, err := apiClient.SendChatCompletion(chatReq)
	if err != nil {
		if apiErr, ok := err.(*api.APIError); ok {
			PrintAPIError(apiErr)
		} else {
			PrintError(err.Error())
		}
		return err
	}

	// Format output
	format := FormatPretty
	if jsonOutput {
		format = FormatJSON
	} else if rawOutput {
		format = FormatRaw
	}

	return FormatChatResponse(resp, format)
}

func init() {
	chatCmd.Flags().StringVarP(&model, "model", "m", "", "Model to use (e.g., openai/gpt-4)")
	chatCmd.Flags().Float64VarP(&temperature, "temperature", "t", 0, "Temperature for response generation (0.0-2.0)")
	chatCmd.Flags().IntVar(&maxTokens, "max-tokens", 0, "Maximum tokens in response")
	chatCmd.Flags().BoolVar(&useStdin, "stdin", false, "Combine argument with piped input (cat file.txt | openrouter chat --stdin 'Analyze:')")
	chatCmd.Flags().BoolVar(&rawOutput, "raw", false, "Output only the response text (no formatting)")
	chatCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output full API response as JSON")
}
