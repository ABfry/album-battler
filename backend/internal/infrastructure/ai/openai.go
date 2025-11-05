package ai

import (
	"context"
	"errors"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

// OpenAIClient implements the Client interface for OpenAI
type OpenAIClient struct {
	client       openai.Client
	defaultModel string
}

// OpenAIConfig holds configuration for OpenAI client
type OpenAIConfig struct {
	APIKey       string
	DefaultModel string // e.g., "gpt-4", "gpt-3.5-turbo", "gpt-4o"
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient(config OpenAIConfig) (*OpenAIClient, error) {
	if config.APIKey == "" {
		return nil, errors.New("OpenAI API key is required")
	}

	if config.DefaultModel == "" {
		config.DefaultModel = "gpt-4"
	}

	client := openai.NewClient(
		option.WithAPIKey(config.APIKey),
	)

	return &OpenAIClient{
		client:       client,
		defaultModel: config.DefaultModel,
	}, nil
}

// Generate generates text using OpenAI's API
func (c *OpenAIClient) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	// Convert messages
	messages := make([]openai.ChatCompletionMessageParamUnion, 0, len(req.Messages))
	for _, msg := range req.Messages {
		switch msg.Role {
		case "user":
			messages = append(messages, openai.UserMessage(msg.Content))
		case "assistant":
			messages = append(messages, openai.AssistantMessage(msg.Content))
		case "system":
			messages = append(messages, openai.SystemMessage(msg.Content))
		default:
			return nil, errors.New("invalid message role: " + msg.Role)
		}
	}

	// Determine model
	model := req.Model
	if model == "" {
		model = c.defaultModel
	}

	// Build request parameters
	params := openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    shared.ChatModel(model),
	}

	if req.MaxTokens > 0 {
		params.MaxTokens = openai.Int(int64(req.MaxTokens))
	}

	if req.Temperature > 0 {
		params.Temperature = openai.Float(float64(req.Temperature))
	}

	// Make API call
	completion, err := c.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(completion.Choices) == 0 {
		return nil, errors.New("no response from OpenAI")
	}

	choice := completion.Choices[0]

	return &GenerateResponse{
		Content:      choice.Message.Content,
		FinishReason: string(choice.FinishReason),
		Usage: TokenUsage{
			PromptTokens:     int(completion.Usage.PromptTokens),
			CompletionTokens: int(completion.Usage.CompletionTokens),
			TotalTokens:      int(completion.Usage.TotalTokens),
		},
	}, nil
}

// Close closes the OpenAI client
func (c *OpenAIClient) Close() error {
	// OpenAI client doesn't require explicit closing
	return nil
}
