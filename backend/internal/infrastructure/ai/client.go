package ai

import "context"

// Message represents a single message in a conversation
type Message struct {
	Role    string // "user", "assistant", "system"
	Content string
}

// GenerateRequest represents a request to generate text
type GenerateRequest struct {
	Messages    []Message
	MaxTokens   int
	Temperature float32
	Model       string // Optional: specific model to use
}

// GenerateResponse represents a response from text generation
type GenerateResponse struct {
	Content      string
	FinishReason string // "stop", "length", "content_filter", etc.
	Usage        TokenUsage
}

// TokenUsage represents token usage information
type TokenUsage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// Client is the interface for AI client implementations
type Client interface {
	// Generate generates text based on the provided messages
	Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)

	// Close closes the client and releases resources
	Close() error
}
