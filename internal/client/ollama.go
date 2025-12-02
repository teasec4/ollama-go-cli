package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// OllamaClient handles communication with Ollama API
type OllamaClient struct {
	baseURL string
	model   string
	client  *http.Client
}

// ChatRequest represents a chat message request to Ollama API
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

// ChatMessage represents a single message in the conversation
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse represents the response from Ollama API
type ChatResponse struct {
	Model     string      `json:"model"`
	CreatedAt string      `json:"created_at"`
	Message   ChatMessage `json:"message"`
	Done      bool        `json:"done"`
	Metrics   Metrics     `json:"metrics,omitempty"`
}

// Metrics contains token usage information
type Metrics struct {
	PromptTokens     int `json:"prompt_tokens"`
	ResponseTokens   int `json:"response_tokens"`
	TotalTokens      int `json:"total_tokens"`
	PromptEvalCount  int `json:"prompt_eval_count"`
	ResponseEvalTime int `json:"response_eval_time"`
}

// NewOllamaClient creates a new Ollama API client
func NewOllamaClient(baseURL, model string) *OllamaClient {
	return &OllamaClient{
		baseURL: baseURL,
		model:   model,
		client:  &http.Client{},
	}
}

// ChatStreamResult holds the final result from streaming
type ChatStreamResult struct {
	Text    string
	Metrics Metrics
}

// ChatStream sends a message to Ollama and streams the response
// Returns a channel for text chunks, a channel for final result, and an error
func (o *OllamaClient) ChatStream(messages []ChatMessage) (<-chan string, <-chan ChatStreamResult, error) {
	req := ChatRequest{
		Model:    o.model,
		Messages: messages,
		Stream:   true,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, nil, err
	}

	httpReq, err := http.NewRequest("POST", o.baseURL+"/api/chat", bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(httpReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	textChan := make(chan string)
	resultChan := make(chan ChatStreamResult, 1)

	// Start goroutine to read streaming response
	go func() {
		defer resp.Body.Close()
		defer close(textChan)

		var finalMetrics Metrics
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Bytes()

			var chatResp ChatResponse
			if err := json.Unmarshal(line, &chatResp); err != nil {
				continue
			}

			// Send text chunk
			if chatResp.Message.Content != "" {
				textChan <- chatResp.Message.Content
			}

			// If done, save metrics
			if chatResp.Done {
				finalMetrics = chatResp.Metrics
			}
		}

		// Send final result
		resultChan <- ChatStreamResult{Metrics: finalMetrics}
		close(resultChan)
	}()

	return textChan, resultChan, nil
}

// Chat sends a message to Ollama and returns the response with token count
func (o *OllamaClient) Chat(messages []ChatMessage) (*ChatResponse, int, error) {
	req := ChatRequest{
		Model:    o.model,
		Messages: messages,
		Stream:   false,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, 0, err
	}

	httpReq, err := http.NewRequest("POST", o.baseURL+"/api/chat", bytes.NewBuffer(body))
	if err != nil {
		return nil, 0, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(httpReq)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, 0, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	// If TotalTokens is 0, calculate from prompt and response counts
	totalTokens := chatResp.Metrics.TotalTokens
	if totalTokens == 0 {
		totalTokens = chatResp.Metrics.PromptTokens + chatResp.Metrics.ResponseTokens
	}

	return &chatResp, totalTokens, nil
}

// SetModel changes the current model
func (o *OllamaClient) SetModel(model string) {
	o.model = model
}

// GetModel returns the current model
func (o *OllamaClient) GetModel() string {
	return o.model
}
