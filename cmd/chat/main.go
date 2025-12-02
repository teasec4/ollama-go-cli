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
	session := chat.NewSession(cfg.Model, cfg.OllamaURL)
	model := ui.NewTUIModel(session)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
