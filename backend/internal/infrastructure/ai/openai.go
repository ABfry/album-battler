package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ABfry/album-battler/backend/internal/domain/service/llm"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

var _ llm.LLMClient = (*OpenAIClient)(nil)

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
func (c *OpenAIClient) Generate(ctx context.Context, req *llm.GenerateRequest) (*llm.GenerateResponse, error) {
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

	return &llm.GenerateResponse{
		Content:      choice.Message.Content,
		FinishReason: string(choice.FinishReason),
		Usage: llm.TokenUsage{
			PromptTokens:     int(completion.Usage.PromptTokens),
			CompletionTokens: int(completion.Usage.CompletionTokens),
			TotalTokens:      int(completion.Usage.TotalTokens),
		},
	}, nil
}

// GenerateThemeは写真撮影バトル用のテーマを生成する
func (c *OpenAIClient) GenerateTheme(ctx context.Context) (string, error) {
	req := &llm.GenerateRequest{
		Messages: []llm.Message{
			{
				Role:    "user",
				Content: GenerateThemePrompt,
			},
		},
		Temperature: 1.0, // 創造性を高めるため高めの温度設定
		MaxTokens:   100, // テーマ名のみなので少なめ
	}

	resp, err := c.Generate(ctx, req)
	if err != nil {
		return "", err
	}

	theme := resp.Content
	if theme == "" {
		return "", errors.New("generated theme is empty")
	}

	return theme, nil
}

// JudgeImageは提出された画像を評価し、スコアと理由を返す
func (c *OpenAIClient) JudgeImage(ctx context.Context, req *llm.ImageJudgeRequest) (*llm.ImageJudgeResponse, error) {
	if req == nil || len(req.ImageURLs) == 0 {
		return nil, errors.New("no images provided")
	}

	results := make([]llm.JudgeResult, 0, len(req.ImageURLs))

	// 各画像を個別に評価
	for _, imageURL := range req.ImageURLs {
		result, err := c.judgeOneImage(ctx, req.Theme, imageURL)
		if err != nil {
			// エラーが発生した画像はスコア0で記録
			results = append(results, llm.JudgeResult{
				Score:  0,
				Reason: fmt.Sprintf("評価エラー: %v", err),
			})
			continue
		}
		results = append(results, result)
	}

	return &llm.ImageJudgeResponse{
		Results: results,
	}, nil
}

// judgeOneImageは1枚の画像を評価する内部ヘルパー関数（Vision API + JSON mode使用）
func (c *OpenAIClient) judgeOneImage(ctx context.Context, theme, imageURL string) (llm.JudgeResult, error) {
	// システムプロンプト（JSON形式での返答を要求）
	systemPrompt := JudgeImagePrompt + "\n必ず以下のJSON形式で返答してください：\n{\"score\": <0-100の整数>, \"reason\": \"<評価理由>\"}"

	// ユーザーメッセージ（テキストと画像の複数パート）
	userPrompt := fmt.Sprintf("お題「%s」に対して、この写真を評価してください。", theme)

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
		openai.UserMessage([]openai.ChatCompletionContentPartUnionParam{
			openai.TextContentPart(userPrompt),
			openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
				URL: imageURL,
			}),
		}),
	}

	// APIリクエストパラメータ
	params := openai.ChatCompletionNewParams{
		Messages:    messages,
		Model:       shared.ChatModel(c.defaultModel),
		Temperature: openai.Float(0.7),
		MaxTokens:   openai.Int(500),
	}

	// APIコール
	completion, err := c.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return llm.JudgeResult{}, fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	if len(completion.Choices) == 0 {
		return llm.JudgeResult{}, errors.New("no response from OpenAI")
	}

	// JSONをパース
	content := completion.Choices[0].Message.Content
	var result struct {
		Score  int    `json:"score"`
		Reason string `json:"reason"`
	}

	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return llm.JudgeResult{}, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return llm.JudgeResult{
		Score:  result.Score,
		Reason: result.Reason,
	}, nil
}

// CloseはOpenAIクライアントをクローズ（明示的な解放は不要）
func (c *OpenAIClient) Close() error {
	// OpenAIクライアントは明示的なクローズ処理が不要
	return nil
}
