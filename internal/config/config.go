package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	ErrNoAPIKey = errors.New("no API key found")
)

// Config represents the application configuration
type Config struct {
	APIKey           string  `yaml:"api_key"`
	DefaultModel     string  `yaml:"default_model"`
	DefaultTemp      float64 `yaml:"default_temperature"`
	DefaultMaxTokens int     `yaml:"default_max_tokens"`
	OutputFormat     string  `yaml:"output_format"`
	APIBaseURL       string  `yaml:"api_base_url"`
	Timeout          int     `yaml:"timeout"`
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		DefaultModel:     "openai/gpt-4",
		DefaultTemp:      1.0,
		DefaultMaxTokens: 4096,
		OutputFormat:     "pretty",
		APIBaseURL:       "https://openrouter.ai/api/v1",
		Timeout:          60,
	}
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	// Try XDG Base Directory spec first
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, "openrouter", "config.yaml")
	}

	// Fall back to ~/.config/openrouter/config.yaml
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fall back to ~/.openrouter.yaml if no home dir
		return filepath.Join(homeDir, ".openrouter.yaml")
	}

	return filepath.Join(homeDir, ".config", "openrouter", "config.yaml")
}

// Load loads configuration from file and environment
func Load() (*Config, error) {
	cfg := DefaultConfig()

	// Try loading from config file
	configPath := GetConfigPath()
	if fileData, err := os.ReadFile(configPath); err == nil {
		if err := yaml.Unmarshal(fileData, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Environment variable takes precedence
	if apiKey := os.Getenv("OPENROUTER_API_KEY"); apiKey != "" {
		cfg.APIKey = apiKey
	}

	// Don't validate API key here - let the command handle it
	// This allows command-line overrides to work properly
	return cfg, nil
}

// Save writes configuration to file
func Save(cfg *Config) error {
	configPath := GetConfigPath()
	configDir := filepath.Dir(configPath)

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file with restricted permissions
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// PartialConfig represents a config that can be missing fields
type PartialConfig struct {
	APIKey           *string  `yaml:"api_key"`
	DefaultModel     *string  `yaml:"default_model"`
	DefaultTemp      *float64 `yaml:"default_temperature"`
	DefaultMaxTokens *int     `yaml:"default_max_tokens"`
	OutputFormat     *string  `yaml:"output_format"`
	APIBaseURL       *string  `yaml:"api_base_url"`
	Timeout          *int     `yaml:"timeout"`
}

// Merge merges a partial config into a full config
func (partial *PartialConfig) Merge(cfg *Config) {
	if partial.APIKey != nil {
		cfg.APIKey = *partial.APIKey
	}
	if partial.DefaultModel != nil {
		cfg.DefaultModel = *partial.DefaultModel
	}
	if partial.DefaultTemp != nil {
		cfg.DefaultTemp = *partial.DefaultTemp
	}
	if partial.DefaultMaxTokens != nil {
		cfg.DefaultMaxTokens = *partial.DefaultMaxTokens
	}
	if partial.OutputFormat != nil {
		cfg.OutputFormat = *partial.OutputFormat
	}
	if partial.APIBaseURL != nil {
		cfg.APIBaseURL = *partial.APIBaseURL
	}
	if partial.Timeout != nil {
		cfg.Timeout = *partial.Timeout
	}
}
