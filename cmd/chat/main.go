package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/teasec4/ollama-go-cli/internal/chat"
	"github.com/teasec4/ollama-go-cli/internal/config"
	"github.com/teasec4/ollama-go-cli/internal/ui"
)

func main() {
	cfg := config.LoadConfig()

	// Create chat session
	session := chat.NewSession("main-session", cfg.Model, cfg.OllamaURL)

	// Create TUI model
	model := ui.NewTUIModel(session)

	// Run TUI
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
