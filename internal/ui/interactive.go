package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/teasec4/ollama-go-cli/internal/chat"
)

// InteractiveSession handles the interactive chat loop
type InteractiveSession struct {
	Session    *chat.Session
	Console    *ConsoleUI
	Input      *bufio.Reader
	MemoryGB   int64
	MemoryUsed int64
}

// NewInteractive creates a new interactive chat session
func NewInteractive(session *chat.Session, console *ConsoleUI) *InteractiveSession {
	is := &InteractiveSession{
		Session: session,
		Console: console,
		Input:   bufio.NewReader(os.Stdin),
		MemoryGB: 120, // Model memory requirement
	}
	is.updateMemoryUsage()
	return is
}

// updateMemoryUsage gets current memory usage from Ollama process
func (s *InteractiveSession) updateMemoryUsage() {
	if runtime.GOOS != "darwin" {
		return
	}

	// Get process memory usage (macOS)
	cmd := exec.Command("ps", "aux")
	output, err := cmd.Output()
	if err != nil {
		return
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "ollama") && !strings.Contains(line, "grep") {
			fields := strings.Fields(line)
			if len(fields) > 5 {
				// RSS is in KB (field 5)
				rss, err := strconv.ParseInt(fields[5], 10, 64)
				if err == nil {
					s.MemoryUsed = rss / 1024 / 1024 // Convert to GB
				}
			}
		}
	}
}

// DrawTopBar displays the session status bar with memory usage
func (s *InteractiveSession) DrawTopBar() {
	s.updateMemoryUsage()

	clearLine := "\033[2K"
	moveToStart := "\033[0G"

	fmt.Print(clearLine + moveToStart)

	// Calculate memory percentage
	memPercent := 0
	if s.MemoryGB > 0 {
		memPercent = int((s.MemoryUsed * 100) / s.MemoryGB)
		if memPercent > 100 {
			memPercent = 100
		}
	}

	// Create memory bar
	barLength := 10
	filledLength := (memPercent * barLength) / 100
	memBar := "["
	for i := 0; i < barLength; i++ {
		if i < filledLength {
			memBar += "="
		} else {
			memBar += " "
		}
	}
	memBar += "]"

	panel := fmt.Sprintf("120k [%d%% %s]", memPercent, memBar)

	width := 80
	padding := width - len(panel)

	if padding < 0 {
		padding = 0
	}

	space := strings.Repeat(" ", padding)

	color.New(color.FgHiYellow).Printf("%s%s\n", space, panel)
}

// ShowThinkingAnimation displays an animated loading indicator
func ShowThinkingAnimation(done <-chan bool) {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0

	for {
		select {
		case <-done:
			fmt.Print("\r")
			return
		default:
			fmt.Printf("\r%s thinking", frames[i%len(frames)])
			i++
			time.Sleep(100 * time.Millisecond)
		}
	}
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

		// Add to session history
		s.Session.AddUserMessage(text)

		// Send request to Ollama with streaming
		fmt.Println() // Newline before Assistant
		color.New(color.FgHiGreen).Print("Assistant: ")

		done := make(chan bool)
		go ShowThinkingAnimation(done)

		messages := s.Session.GetMessages()
		textChan, resultChan, err := s.Session.Client.ChatStream(messages)

		if err != nil {
			done <- true
			color.New(color.FgHiRed).Printf("Error: %v\n\n", err)
			// Remove last user message on error
			s.Session.Messages = s.Session.Messages[:len(s.Session.Messages)-1]
			continue
		}

		done <- true

		// Stream response text
		var fullResponse strings.Builder
		for chunk := range textChan {
			fmt.Print(chunk)
			fullResponse.WriteString(chunk)
		}

		// Get final result
		result := <-resultChan

		// Add assistant message to session
		responseText := fullResponse.String()
		s.Session.AddAssistantMessage(responseText)

		// Calculate and add tokens
		totalTokens := result.Metrics.TotalTokens
		if totalTokens == 0 {
			totalTokens = result.Metrics.PromptTokens + result.Metrics.ResponseTokens
		}

		// If API doesn't return tokens, estimate them
		if totalTokens == 0 {
			// Get the last two messages (user + assistant)
			if len(s.Session.Messages) >= 2 {
				userMsg := s.Session.Messages[len(s.Session.Messages)-2].Content
				assistantMsg := s.Session.Messages[len(s.Session.Messages)-1].Content
				totalTokens = chat.EstimateTokens(userMsg) + chat.EstimateTokens(assistantMsg)
			}
		}

		s.Session.AddTokens(totalTokens)

		fmt.Println()
	}
}
