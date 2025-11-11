package ai

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/genai"
)

// Gemini向けのClientインターフェース実装
type GeminiClient struct {
	client       *genai.Client
	defaultModel string
}

// GeminiConfigはGeminiクライアントの設定を保持
type GeminiConfig struct {
	APIKey       string
	DefaultModel string // 例: "gemini-2.0-flash", "gemini-1.5-pro"
}

// NewGeminiClientは新しいGeminiクライアントを作成
func NewGeminiClient(ctx context.Context, config GeminiConfig) (*GeminiClient, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("gemini API key is required")
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

// GenerateはGeminiのAPIでテキスト生成を行う
func (c *GeminiClient) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	if len(req.Messages) == 0 {
		return nil, errors.New("no messages provided")
	}

	// モデルを決定
	modelName := req.Model
	if modelName == "" {
		modelName = c.defaultModel
	}

	// 生成設定を構築
	config := &genai.GenerateContentConfig{}

	if req.MaxTokens > 0 {
		config.MaxOutputTokens = int32(req.MaxTokens)
	}

	if req.Temperature > 0 {
		config.Temperature = genai.Ptr(req.Temperature)
	}

	// メッセージをGemini形式に変換
	var history []*genai.Content
	var systemInstruction string
	messages := req.Messages

	// 先頭がsystemだった場合システムメッセージとして抽出
	if len(messages) > 0 && messages[0].Role == "system" {
		systemInstruction = messages[0].Content
		messages = messages[1:]

		config.SystemInstruction = &genai.Content{
			Parts: []*genai.Part{
				{Text: systemInstruction},
			},
		}
	}

	// 会話履歴を構築（最後のユーザーメッセージは除外）
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

	// 最後のメッセージがユーザーか検証
	lastMsg := messages[len(messages)-1]
	if lastMsg.Role != "user" {
		return nil, errors.New("last message must be from user")
	}

	// 履歴付きでチャットセッション作成
	chat, err := c.client.Chats.Create(ctx, modelName, config, history)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat session: %w", err)
	}

	// 最後のメッセージを送信
	resp, err := chat.SendMessage(ctx, genai.Part{Text: lastMsg.Content})
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	// テキスト内容を抽出
	content := resp.Text()

	// 終了理由をマッピング
	finishReason := "stop"
	if len(resp.Candidates) > 0 {
		finishReason = fmt.Sprintf("%v", resp.Candidates[0].FinishReason)
	}

	// トークン使用量を抽出
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

// CloseはGeminiクライアントのクローズ処理（明示的なクローズは不要）
func (c *GeminiClient) Close() error {
	// 新しいgenaiクライアントは明示的なクローズ不要
	return nil
}
