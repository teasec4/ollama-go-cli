package client

import (
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

	totalTokens := chatResp.Metrics.TotalTokens
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
