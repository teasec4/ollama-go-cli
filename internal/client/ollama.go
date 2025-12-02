package client

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

// OllamaClient handles communication with Ollama API
type OllamaClient struct {
	client *openai.Client
	model  string
}

// ChatMessage represents a single message in the conversation
type ChatMessage = openai.ChatCompletionMessage

// Metrics contains token usage information
type Metrics struct {
	PromptTokens   int
	ResponseTokens int
	TotalTokens    int
}

// ChatStreamResult holds the final result from streaming
type ChatStreamResult struct {
	Text    string
	Metrics Metrics
}

// NewOllamaClient creates a new Ollama API client
// baseURL should be like "http://localhost:11434/v1"
func NewOllamaClient(baseURL, model string) *OllamaClient {
	config := openai.DefaultConfig("ollama") // dummy key, unused by Ollama
	config.BaseURL = baseURL
	return &OllamaClient{
		client: openai.NewClientWithConfig(config),
		model:  model,
	}
}

// ChatStream sends a message to Ollama and streams the response
// Returns a channel for text chunks, a channel for final result, and an error
func (o *OllamaClient) ChatStream(messages []ChatMessage) (<-chan string, <-chan ChatStreamResult, error) {
	textChan := make(chan string)
	resultChan := make(chan ChatStreamResult, 1)

	req := openai.ChatCompletionRequest{
		Model:    o.model,
		Messages: messages,
		Stream:   true,
	}

	// Start goroutine to read streaming response
	go func() {
		defer close(textChan)
		defer close(resultChan)

		stream, err := o.client.CreateChatCompletionStream(context.Background(), req)
		if err != nil {
			resultChan <- ChatStreamResult{Metrics: Metrics{}}
			return
		}
		defer stream.Close()

		var fullText string
		for {
			response, err := stream.Recv()
			if err != nil {
				break
			}

			// Send text chunk
			if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
				chunk := response.Choices[0].Delta.Content
				textChan <- chunk
				fullText += chunk
			}
		}

		// Send final result with estimated tokens
		resultChan <- ChatStreamResult{
			Text:    fullText,
			Metrics: Metrics{TotalTokens: 0}, // Ollama may not return token counts in stream
		}
	}()

	return textChan, resultChan, nil
}

// Chat sends a message to Ollama and returns the response with token count
func (o *OllamaClient) Chat(messages []ChatMessage) (string, int, error) {
	req := openai.ChatCompletionRequest{
		Model:    o.model,
		Messages: messages,
		Stream:   false,
	}

	resp, err := o.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", 0, fmt.Errorf("no choices returned from API")
	}

	text := resp.Choices[0].Message.Content
	totalTokens := resp.Usage.TotalTokens

	return text, totalTokens, nil
}

// SetModel changes the current model
func (o *OllamaClient) SetModel(model string) {
	o.model = model
}

// GetModel returns the current model
func (o *OllamaClient) GetModel() string {
	return o.model
}
