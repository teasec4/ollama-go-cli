package config

// Config holds application configuration
type Config struct {
	OllamaURL     string
	Model         string
	AssistantName string
	TerminalWidth int
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		OllamaURL:     "http://localhost:11434",
		Model:         "gpt-oss:20b",
		AssistantName: "Max",
		TerminalWidth: 80,
	}
}

// LoadConfig loads the configuration
// TODO: implement loading from environment variables or config file
func LoadConfig() *Config {
	return DefaultConfig()
}
