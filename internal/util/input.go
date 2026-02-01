package util

import (
	"errors"
	"io"
	"os"
	"strings"
)

// GetInput reads input from either command-line arguments or stdin
func GetInput(args []string) (string, error) {
	// If arguments provided, use them
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	// Check if there's input from stdin (pipe)
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", errors.New("failed to check stdin")
	}

	// Check if stdin is a pipe (not a terminal character device)
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Reading from pipe
		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", errors.New("failed to read from stdin")
		}
		return strings.TrimSpace(string(bytes)), nil
	}

	// No input provided
	return "", errors.New("no input provided: use argument or pipe")
}

// CombineInputWithStdin combines command-line arguments with stdin content
// If useStdin is true, appends stdin to the prompt argument
// Useful for: cat file.txt | openrouter chat --stdin "Analyze this:"
func CombineInputWithStdin(args []string, useStdin bool) (string, error) {
	if !useStdin {
		// Fallback to normal input handling
		return GetInput(args)
	}

	// Get the prompt from arguments
	var prompt string
	if len(args) > 0 {
		prompt = strings.Join(args, " ")
	}

	// Check if stdin has data
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", errors.New("failed to check stdin")
	}

	// Check if stdin is a pipe
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Read stdin
		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", errors.New("failed to read from stdin")
		}
		stdinContent := strings.TrimSpace(string(bytes))

		// Combine prompt and stdin
		if prompt != "" {
			return prompt + "\n\n" + stdinContent, nil
		}
		return stdinContent, nil
	}

	// No stdin, just use prompt
	if prompt != "" {
		return prompt, nil
	}

	return "", errors.New("no input provided: use argument or pipe")
}
