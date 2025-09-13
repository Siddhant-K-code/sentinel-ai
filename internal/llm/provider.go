package llm

import (
	"context"
	"encoding/json"
)

// Model represents an LLM interface
type Model interface {
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
}

// ChatRequest represents a chat request
type ChatRequest struct {
	System     string    `json:"system"`
	Messages   []Message `json:"messages"`
	Tools      []Tool    `json:"tools,omitempty"`
	MaxTokens  int       `json:"max_tokens"`
	Temperature float32  `json:"temperature"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Tool represents a tool definition
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  interface{} `json:"parameters"`
}

// ChatResponse represents a chat response
type ChatResponse struct {
	Content   string      `json:"content"`
	ToolCalls []ToolCall  `json:"tool_calls,omitempty"`
	Usage     Usage       `json:"usage,omitempty"`
}

// ToolCall represents a tool call
type ToolCall struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"`
	Function ToolCallFunction `json:"function"`
}

// ToolCallFunction represents a tool call function
type ToolCallFunction struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// Usage represents token usage
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Provider represents an LLM provider
type Provider interface {
	GetModel(alias string) (Model, error)
	ListModels() []string
}

// Config represents LLM configuration
type Config struct {
	PrimaryAlias   string
	SecondaryAlias string
	APIKey         string
	BaseURL        string
	Timeout        int
}

// NewProvider creates a new LLM provider
func NewProvider(config Config) (Provider, error) {
	// TODO: Implement provider selection based on config
	// This would support OpenAI, Anthropic, etc.
	return &mockProvider{}, nil
}

// mockProvider is a placeholder implementation
type mockProvider struct{}

func (m *mockProvider) GetModel(alias string) (Model, error) {
	return &mockModel{}, nil
}

func (m *mockProvider) ListModels() []string {
	return []string{"gpt-4", "claude-3-sonnet"}
}

// mockModel is a placeholder implementation
type mockModel struct{}

func (m *mockModel) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	// TODO: Implement actual LLM calls
	return &ChatResponse{
		Content: "Mock response",
		Usage: Usage{
			PromptTokens:     100,
			CompletionTokens: 50,
			TotalTokens:      150,
		},
	}, nil
}
