package ui

import (
	"fmt"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/teasec4/ollama-go-cli/internal/chat"
)

// TUIModel is the main TUI model
type TUIModel struct {
	session      *chat.Session
	input        string
	cursorPos    int
	width        int
	height       int
	thinking     bool
	memPercent   int
	err          error
	currentReply string
	thinkFrame   int
	scrollOffset int // For scrolling through message history
	mu           sync.Mutex // Protect concurrent access
}

// NewTUIModel creates a new TUI model
func NewTUIModel(session *chat.Session) *TUIModel {
	return &TUIModel{
		session:    session,
		input:      "",
		width:      80,
		height:     24,
		thinking:   false,
		memPercent: 0,
		thinkFrame: 0,
	}
}

// Init initializes the model
func (m *TUIModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		if m.thinking {
			// Allow scrolling even while thinking
			switch msg.String() {
			case "pageup", "up":
				m.scrollOffset += 3
				return m, nil
			case "pagedown", "down":
				m.scrollOffset -= 3
				if m.scrollOffset < 0 {
					m.scrollOffset = 0
				}
				return m, nil
			}
			return m, nil // Don't accept other input while thinking
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "pageup":
			m.scrollOffset += 5
			return m, nil
		case "pagedown":
			m.scrollOffset -= 5
			if m.scrollOffset < 0 {
				m.scrollOffset = 0
			}
			return m, nil
		case "up":
			m.scrollOffset += 1
			return m, nil
		case "down":
			m.scrollOffset -= 1
			if m.scrollOffset < 0 {
				m.scrollOffset = 0
			}
			return m, nil
		case "enter":
			if m.input != "" {
				userMsg := m.input
				m.input = ""
				m.cursorPos = 0
				m.scrollOffset = 0 // Reset scroll when sending message

				// Add to session
				m.session.AddUserMessage(userMsg)

				// Send message
				m.thinking = true
				m.currentReply = ""
				m.thinkFrame = 0
				return m, tea.Batch(m.sendMessageCmd(), m.animateThinker())
			}
		case "backspace":
			if m.cursorPos > 0 {
				runes := []rune(m.input)
				m.input = string(append(runes[:m.cursorPos-1], runes[m.cursorPos:]...))
				m.cursorPos--
			}
		case "left":
			if m.cursorPos > 0 {
				m.cursorPos--
			}
		case "right":
			if m.cursorPos < len([]rune(m.input)) {
				m.cursorPos++
			}
		case "home":
			m.cursorPos = 0
		case "end":
			m.cursorPos = len([]rune(m.input))
		default:
			runes := []rune(m.input)
			runes = append(runes[:m.cursorPos], append([]rune(msg.String()), runes[m.cursorPos:]...)...)
			m.input = string(runes)
			m.cursorPos += len([]rune(msg.String()))
		}

	case msgChunk:
		m.mu.Lock()
		m.currentReply += msg.text
		m.mu.Unlock()
		return m, nil

	case msgThinkingDone:
		m.mu.Lock()
		m.thinking = false
		m.session.AddAssistantMessage(m.currentReply)
		tokensUsed := chat.EstimateTokens(m.currentReply)
		m.session.AddTokens(tokensUsed)
		m.currentReply = ""
		m.mu.Unlock()

	case msgThinkFrame:
		if m.thinking {
			m.thinkFrame = (m.thinkFrame + 1) % 10
			return m, m.animateThinker()
		}

	case msgError:
		m.mu.Lock()
		m.thinking = false
		m.err = msg.err
		if len(m.session.Messages) > 0 {
			m.session.Messages = m.session.Messages[:len(m.session.Messages)-1]
		}
		m.mu.Unlock()
	}

	return m, nil
}

// View renders the UI
func (m *TUIModel) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	// Styles
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Padding(0, 1)

	msgWidth := m.width - 3
	if msgWidth < 20 {
		msgWidth = 20
	}

	assistantStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")).
		Bold(true).
		PaddingLeft(1)

	userStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("51")).
		PaddingLeft(1)

	assistantLabelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("40")).
		Bold(true)

	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("51")).
		Bold(true).
		PaddingLeft(1)

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")).
		PaddingRight(1).
		PaddingLeft(1)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	// Header
	header := headerStyle.Render("ollama-chat (local)")

	// Messages area
	messagesHeight := m.height - 6
	if messagesHeight < 1 {
		messagesHeight = 1
	}

	var messageLines []string

	m.mu.Lock()
	// Add session history
	for _, msg := range m.session.Messages {
		if msg.Role == "user" {
			text := fmt.Sprintf("You: %s", msg.Content)
			wrappedText := wrapText(text, msgWidth)
			for _, line := range strings.Split(wrappedText, "\n") {
				if line != "" {
					messageLines = append(messageLines, userStyle.Render(line))
				}
			}
		} else {
			text := fmt.Sprintf("Assistant: %s", msg.Content)
			wrappedText := wrapText(text, msgWidth)
			lines := strings.Split(wrappedText, "\n")
			for i, line := range lines {
				if line != "" {
					if i == 0 {
						// First line: bold label + content
						parts := strings.SplitN(line, ": ", 2)
						if len(parts) == 2 {
							rendered := assistantLabelStyle.Render(parts[0]+": ") + assistantStyle.Render(parts[1])
							messageLines = append(messageLines, rendered)
						} else {
							messageLines = append(messageLines, assistantStyle.Render(line))
						}
					} else {
						messageLines = append(messageLines, assistantStyle.Render(line))
					}
				}
			}
		}
	}

	// Add current reply if streaming
	if m.thinking && m.currentReply != "" {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		replyDisplay := m.currentReply + frames[m.thinkFrame]
		wrappedReply := wrapText(replyDisplay, msgWidth)
		replyLines := strings.Split(wrappedReply, "\n")
		for i, line := range replyLines {
			if i == 0 {
				rendered := assistantLabelStyle.Render("Assistant: ") + assistantStyle.Render(line)
				messageLines = append(messageLines, rendered)
			} else {
				messageLines = append(messageLines, assistantStyle.Render(line))
			}
		}
	}
	m.mu.Unlock()

	// Apply scroll offset
	if m.scrollOffset > 0 {
		if m.scrollOffset >= len(messageLines) {
			messageLines = []string{}
		} else {
			messageLines = messageLines[:len(messageLines)-m.scrollOffset]
		}
	}

	// Keep last N messages
	if len(messageLines) > messagesHeight {
		messageLines = messageLines[len(messageLines)-messagesHeight:]
	}

	messagesView := strings.Join(messageLines, "\n")
	if len(messageLines) < messagesHeight {
		messagesView += strings.Repeat("\n", messagesHeight-len(messageLines))
	}

	// Input area
	inputPrompt := inputStyle.Render("You: ")
	var inputArea string

	inputMaxWidth := m.width - 6
	runes := []rune(m.input)
	var displayInput string

	if m.thinking {
		displayInput = m.input
	} else {
		displayInput = string(runes[:m.cursorPos]) + "█" + string(runes[m.cursorPos:])
	}

	displayRunes := []rune(displayInput)
	if len(displayRunes) > inputMaxWidth {
		displayInput = "…" + string(displayRunes[len(displayRunes)-inputMaxWidth+1:])
	}

	inputArea = inputPrompt + displayInput

	// Thinking indicator
	thinkingMsg := ""
	if m.thinking && m.currentReply == "" {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		thinkingMsg = assistantLabelStyle.Render("Assistant: ") + assistantStyle.Render(fmt.Sprintf("%s thinking...", frames[m.thinkFrame])) + "\n"
	}

	// Error message
	errorMsg := ""
	if m.err != nil {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
		errorMsg = errorStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n"
		m.err = nil
	}

	// Memory bar
	barLength := 10
	filledLength := (m.memPercent * barLength) / 100
	memBar := "["
	for i := 0; i < barLength; i++ {
		if i < filledLength {
			memBar += "="
		} else {
			memBar += " "
		}
	}
	memBar += "]"

	scrollInfo := ""
	if m.scrollOffset > 0 {
		scrollInfo = fmt.Sprintf(" | ↑ %d lines", m.scrollOffset)
	}
	statusBar := statusStyle.Render(fmt.Sprintf("120k [%d%% %s]%s", m.memPercent, memBar, scrollInfo))

	// Help text
	helpContent := "PgUp/↑:scroll up  PgDn/↓:scroll down  Ctrl+C:quit"
	helpText := helpStyle.Render(helpContent)
	
	// Right-align help text (use visible length for calculation)
	helpLen := visibleLen(helpText)
	helpPadding := m.width - helpLen
	if helpPadding > 0 {
		helpText = strings.Repeat(" ", helpPadding) + helpText
	}

	// Build final view
	view := fmt.Sprintf(
		"%s\n%s\n%s%s%s%s\n%s\n%s\n%s",
		header,
		strings.Repeat("─", m.width),
		messagesView,
		thinkingMsg,
		errorMsg,
		inputArea,
		strings.Repeat("─", m.width),
		statusBar,
		helpText,
	)

	return view
}

// Message types
type msgChunk struct {
	text string
}

type msgThinkingDone struct{}

type msgThinkFrame struct{}

type msgError struct {
	err error
}

// visibleLen returns visible length excluding ANSI codes
func visibleLen(text string) int {
	visible := text
	for {
		start := strings.Index(visible, "\x1b[")
		if start == -1 {
			break
		}
		end := strings.Index(visible[start:], "m")
		if end == -1 {
			break
		}
		visible = visible[:start] + visible[start+end+1:]
	}
	return len(visible)
}

// wrapText wraps text to fit within max width
func wrapText(text string, maxWidth int) string {
	if maxWidth < 1 {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	var currentLine string

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		if visibleLen(testLine) <= maxWidth {
			currentLine = testLine
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}

// animateThinker sends animation frames
func (m *TUIModel) animateThinker() tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg {
		return msgThinkFrame{}
	})
}

// sendMessageCmd sends message to Ollama and updates UI
func (m *TUIModel) sendMessageCmd() tea.Cmd {
	return func() tea.Msg {
		messages := m.session.GetMessages()
		textChan, resultChan, err := m.session.Client.ChatStream(messages)

		if err != nil {
			return msgError{err}
		}

		// Stream chunks directly into model
		for chunk := range textChan {
			m.mu.Lock()
			m.currentReply += chunk
			m.mu.Unlock()
		}

		// Wait for metrics
		<-resultChan

		return msgThinkingDone{}
	}
}
