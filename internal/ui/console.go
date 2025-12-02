package ui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// ConsoleUI provides colored console output for the chat application
type ConsoleUI struct {
	assistantName string
}

// NewConsoleUI creates a new console UI renderer
func NewConsoleUI(assistantName string) *ConsoleUI {
	return &ConsoleUI{assistantName: assistantName}
}

// RenderHeader prints the application header
func (c *ConsoleUI) RenderHeader() {
	title := "=== ollama-chat (local) ==="
	bold := color.New(color.FgHiMagenta, color.Bold)
	bold.Println(title)

	nameStyle := color.New(color.FgHiBlue).Add(color.Bold)
	nameStyle.Printf("Assistant: %s\n", c.assistantName)

	fmt.Println(strings.Repeat("─", 40))
}
