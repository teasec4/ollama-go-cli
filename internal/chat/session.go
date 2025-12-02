package chat

import (
	openai "github.com/sashabaranov/go-openai"
	"github.com/teasec4/ollama-go-cli/internal/client"
)

// Session manages chat state: history and client connection
type Session struct {
	Messages       []openai.ChatCompletionMessage
	TokenCount     int
	Client         *client.OllamaClient
	ContextMaxSize int // Maximum context window size
}

// NewSession creates a new chat session
func NewSession(model, ollamaURL string) *Session {
	return &Session{
		Messages:       []openai.ChatCompletionMessage{},
		TokenCount:     0,
		Client:         client.NewOllamaClient(ollamaURL, model),
		ContextMaxSize: 4096, // Default context size, can be updated by client
	}
}

// AddUserMessage adds a user message to the session history
func (s *Session) AddUserMessage(text string) {
	s.Messages = append(s.Messages, openai.ChatCompletionMessage{
		Role:    "user",
		Content: text,
	})
}

// AddAssistantMessage adds an assistant message to the session history
func (s *Session) AddAssistantMessage(text string) {
	s.Messages = append(s.Messages, openai.ChatCompletionMessage{
		Role:    "assistant",
		Content: text,
	})
}

// AddTokens increments the total token count
func (s *Session) AddTokens(count int) {
	s.TokenCount += count
}

// Clear resets the session
func (s *Session) Clear() {
	s.Messages = []openai.ChatCompletionMessage{}
	s.TokenCount = 0
}
