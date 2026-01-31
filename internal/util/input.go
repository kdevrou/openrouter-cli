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
