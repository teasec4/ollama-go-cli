package ui

import (
	"github.com/fatih/color"
)

// ConsoleUI provides colored console output for the chat application
type ConsoleUI struct{}

// NewConsoleUI creates a new console UI renderer
func NewConsoleUI() *ConsoleUI {
	RenderHeader()
	return &ConsoleUI{}
}

// RenderHeader prints the application header
func RenderHeader() {
	title := "-------------Ollama-chat (local)-------------"
	bold := color.New(color.FgHiMagenta, color.Bold)
	bold.Println(title)
}
