package ai

import (
	"context"
	"errors"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

var _ Client = (*OpenAIClient)(nil)

// OpenAIClientはOpenAI用のClientインターフェース実装
type OpenAIClient struct {
	client       openai.Client
	defaultModel string
}

// OpenAIConfigはOpenAIクライアントの設定を保持
type OpenAIConfig struct {
	APIKey       string
	DefaultModel string // 例: "gpt-4", "gpt-3.5-turbo", "gpt-4o"
}

// NewOpenAIClientは新しいOpenAIクライアントを作成
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

// GenerateはOpenAIのAPIでテキスト生成を行う
func (c *OpenAIClient) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	// メッセージをOpenAI形式に変換
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

	// モデルを決定
	model := req.Model
	if model == "" {
		model = c.defaultModel
	}

	// リクエストパラメータを構築
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

	// APIコール実施
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

func (c *OpenAIClient) GenerateTheme(ctx context.Context) (string, error) {
	return "", nil
}

func (c *OpenAIClient) JudgeImage(ctx context.Context, req *ImageJudgeRequest) (*ImageJudgeResponse, error) {
	return nil, nil
}

// CloseはOpenAIクライアントをクローズ（明示的な解放は不要）
func (c *OpenAIClient) Close() error {
	// OpenAIクライアントは明示的なクローズ処理が不要
	return nil
}
