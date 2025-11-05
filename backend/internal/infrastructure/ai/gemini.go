package ai

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/genai"
)

// GeminiClient implements the Client interface for Google Gemini
type GeminiClient struct {
	client       *genai.Client
	defaultModel string
}

// GeminiConfig holds configuration for Gemini client
type GeminiConfig struct {
	APIKey       string
	DefaultModel string // e.g., "gemini-2.0-flash", "gemini-1.5-pro"
}

// NewGeminiClient creates a new Gemini client
func NewGeminiClient(ctx context.Context, config GeminiConfig) (*GeminiClient, error) {
	if config.APIKey == "" {
		return nil, errors.New("Gemini API key is required")
	}

	if config.DefaultModel == "" {
		config.DefaultModel = "gemini-2.0-flash"
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &GeminiClient{
		client:       client,
		defaultModel: config.DefaultModel,
	}, nil
}

// Generate generates text using Gemini's API
func (c *GeminiClient) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	if len(req.Messages) == 0 {
		return nil, errors.New("no messages provided")
	}

	// Determine model
	modelName := req.Model
	if modelName == "" {
		modelName = c.defaultModel
	}

	// Build generation config
	config := &genai.GenerateContentConfig{}

	if req.MaxTokens > 0 {
		config.MaxOutputTokens = int32(req.MaxTokens)
	}

	if req.Temperature > 0 {
		config.Temperature = genai.Ptr(req.Temperature)
	}

	// Convert messages to Gemini format
	var history []*genai.Content
	var systemInstruction string
	messages := req.Messages

	// Extract system message if present
	if len(messages) > 0 && messages[0].Role == "system" {
		systemInstruction = messages[0].Content
		messages = messages[1:]

		config.SystemInstruction = &genai.Content{
			Parts: []*genai.Part{
				{Text: systemInstruction},
			},
		}
	}

	// Build conversation history (exclude the last user message)
	for i := 0; i < len(messages)-1; i++ {
		msg := messages[i]
		var role string
		switch msg.Role {
		case "user":
			role = "user"
		case "assistant":
			role = "model"
		default:
			return nil, errors.New("invalid message role: " + msg.Role)
		}

		history = append(history, &genai.Content{
			Role: role,
			Parts: []*genai.Part{
				{Text: msg.Content},
			},
		})
	}

	// Validate last message is from user
	lastMsg := messages[len(messages)-1]
	if lastMsg.Role != "user" {
		return nil, errors.New("last message must be from user")
	}

	// Create chat session with history
	chat, err := c.client.Chats.Create(ctx, modelName, config, history)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat session: %w", err)
	}

	// Send the last message
	resp, err := chat.SendMessage(ctx, genai.Part{Text: lastMsg.Content})
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	// Extract text content
	content := resp.Text()

	// Map finish reason
	finishReason := "stop"
	if len(resp.Candidates) > 0 {
		finishReason = fmt.Sprintf("%v", resp.Candidates[0].FinishReason)
	}

	// Extract token usage
	usage := TokenUsage{}
	if resp.UsageMetadata != nil {
		usage.PromptTokens = int(resp.UsageMetadata.PromptTokenCount)
		usage.CompletionTokens = int(resp.UsageMetadata.CandidatesTokenCount)
		usage.TotalTokens = int(resp.UsageMetadata.TotalTokenCount)
	}

	return &GenerateResponse{
		Content:      content,
		FinishReason: finishReason,
		Usage:        usage,
	}, nil
}

// Close closes the Gemini client
func (c *GeminiClient) Close() error {
	// The new genai client doesn't require explicit closing
	return nil
}
