package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/teasec4/ollama-go-cli/internal/chat"
)

// InteractiveSession handles the interactive chat loop
type InteractiveSession struct {
	Session *chat.Session
	Console *ConsoleUI
	Input   *bufio.Reader
}

// NewInteractive creates a new interactive chat session
func NewInteractive(session *chat.Session, console *ConsoleUI) *InteractiveSession {
	return &InteractiveSession{
		Session: session,
		Console: console,
		Input:   bufio.NewReader(os.Stdin),
	}
}

// DrawTopBar displays the session status bar at the top right
func (s *InteractiveSession) DrawTopBar() {
	clearLine := "\033[2K"
	moveToStart := "\033[0G"

	fmt.Print(clearLine + moveToStart)

	panel := fmt.Sprintf("[%s] [Tokens: %d]", s.Session.Name, s.Session.TokenCount)

	width := 80
	padding := width - len(panel)

	if padding < 0 {
		padding = 0
	}

	space := strings.Repeat(" ", padding)

	color.New(color.FgHiYellow).Printf("%s%s\n", space, panel)
}

// Loop starts the main interactive chat loop
func (s *InteractiveSession) Loop() {
	for {
		s.DrawTopBar()

		// Read user input
		userPrompt := color.New(color.FgHiCyan, color.Bold)
		userPrompt.Print("You: ")

		text, _ := s.Input.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "" {
			continue
		}

		if text == "exit" || text == "quit" {
			fmt.Println("👋 Goodbye…")
			return
		}

		// Display user message
		color.New(color.FgHiCyan).Printf("You: %s\n", text)

		// Add to session history
		s.Session.AddUserMessage(text)

		// Send request to Ollama
		color.New(color.FgHiGreen).Print("Max: ")
		fmt.Print("🤔 thinking")

		messages := s.Session.GetMessages()
		resp, tokens, err := s.Session.Client.Chat(messages)

		fmt.Print("\r")

		if err != nil {
			color.New(color.FgHiRed).Printf("Error: %v\n\n", err)
			// Remove last user message on error
			s.Session.Messages = s.Session.Messages[:len(s.Session.Messages)-1]
			continue
		}

		// Display assistant response
		reply := resp.Message.Content
		s.Session.AddAssistantMessage(reply)
		s.Session.AddTokens(tokens)

		color.New(color.FgHiGreen).Printf("Max: %s\n\n", reply)
	}
}
