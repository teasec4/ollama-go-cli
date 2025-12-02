package chat

import (
	"github.com/teasec4/ollama-go-cli/internal/client"
	"github.com/teasec4/ollama-go-cli/internal/constants"
)

// Session manages chat state: history, token count, and client connection
type Session struct {
	Name       string
	Model      string
	Messages   []Message
	TokenCount int
	Client     *client.OllamaClient
}

// Message represents a single message in the conversation
type Message struct {
	Role    string // "user" or "assistant"
	Content string
}

// NewSession creates a new chat session
func NewSession(name, model, ollamaURL string) *Session {
	return &Session{
		Name:       name,
		Model:      model,
		Messages:   []Message{},
		TokenCount: 0,
		Client:     client.NewOllamaClient(ollamaURL, model),
	}
}

// AddUserMessage adds a user message to the session history
func (s *Session) AddUserMessage(text string) {
	s.Messages = append(s.Messages, Message{
		Role:    constants.RoleUser,
		Content: text,
	})
}

// AddAssistantMessage adds an assistant message to the session history
func (s *Session) AddAssistantMessage(text string) {
	s.Messages = append(s.Messages, Message{
		Role:    constants.RoleAssistant,
		Content: text,
	})
}

// AddTokens increments the total token count
func (s *Session) AddTokens(count int) {
	s.TokenCount += count
}

// GetMessages converts session messages to API format
func (s *Session) GetMessages() []client.ChatMessage {
	var msgs []client.ChatMessage
	for _, m := range s.Messages {
		msgs = append(msgs, client.ChatMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}
	return msgs
}

// Clear resets the session (clears history and token count)
func (s *Session) Clear() {
	s.Messages = []Message{}
	s.TokenCount = 0
}
