package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	OllamaURL string
	Model     string
}

// LoadConfig loads the configuration from environment variables or defaults
func LoadConfig() *Config {
	cfg := &Config{
		OllamaURL: "http://localhost:11434/v1",
		Model:     "llama3:latest",
	}

	// Override from environment variables if set
	if url := os.Getenv("OLLAMA_URL"); url != "" {
		cfg.OllamaURL = url
	}
	if model := os.Getenv("OLLAMA_MODEL"); model != "" {
		cfg.Model = model
	}

	return cfg
}
