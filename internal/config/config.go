package config

import (
	"os"

	"github.com/teasec4/ollama-go-cli/internal/constants"
)

// Config holds application configuration
type Config struct {
	OllamaURL     string
	Model         string
	TerminalWidth int
	MemoryGB      int64
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		OllamaURL:     "http://localhost:11434",
		Model:         "llama3:latest",
		TerminalWidth: constants.DefaultTerminalWidth,
		MemoryGB:      constants.DefaultMemoryGB,
	}
}

// LoadConfig loads the configuration from environment variables or defaults
func LoadConfig() *Config {
	cfg := DefaultConfig()

	// Load from environment variables if set
	if url := os.Getenv("OLLAMA_URL"); url != "" {
		cfg.OllamaURL = url
	}
	if model := os.Getenv("OLLAMA_MODEL"); model != "" {
		cfg.Model = model
	}

	return cfg
}
