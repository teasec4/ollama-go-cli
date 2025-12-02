package main

import (
	"github.com/teasec4/ollama-go-cli/internal/chat"
	"github.com/teasec4/ollama-go-cli/internal/config"
	"github.com/teasec4/ollama-go-cli/internal/ui"
)

func main() {
	cfg := config.LoadConfig()

	// UI
	console := ui.NewConsoleUI(cfg.AssistantName)
	console.RenderHeader()

	// Session
	session := chat.NewSession("main-session", cfg.Model, cfg.OllamaURL)

	// Interactive loop
	interactiveUI := ui.NewInteractive(session, console)
	interactiveUI.Loop()
}
