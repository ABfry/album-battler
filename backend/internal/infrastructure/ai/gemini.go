package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

var _ Client = (*GeminiClient)(nil)

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

// GenerateThemeは写真撮影バトル用のテーマを生成する
func (c *GeminiClient) GenerateTheme(ctx context.Context) (string, error) {
	req := &GenerateRequest{
		Messages: []Message{
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
		return "", fmt.Errorf("failed to generate theme: %w", err)
	}

	// 生成されたテーマから余分な空白や改行を除去
	theme := resp.Content
	if theme == "" {
		return "", errors.New("generated theme is empty")
	}

	return theme, nil
}

// JudgeImageは提出された画像を評価し、スコアと理由を返す
func (c *GeminiClient) JudgeImage(ctx context.Context, req *ImageJudgeRequest) (*ImageJudgeResponse, error) {
	if req == nil || len(req.ImageURLs) == 0 {
		return nil, errors.New("no images provided")
	}

	results := make([]JudgeResult, 0, len(req.ImageURLs))

	// 各画像を個別に評価
	for _, imageURL := range req.ImageURLs {
		result, err := c.judgeOneImage(ctx, imageURL)
		if err != nil {
			// エラーが発生した画像はスコア0で記録
			results = append(results, JudgeResult{
				Score:  0,
				Reason: fmt.Sprintf("評価エラー: %v", err),
			})
			continue
		}
		results = append(results, result)
	}

	return &ImageJudgeResponse{
		Results: results,
	}, nil
}

// judgeOneImageは1枚の画像を評価する内部ヘルパー関数
func (c *GeminiClient) judgeOneImage(ctx context.Context, imageURL string) (JudgeResult, error) {
	// チャット生成設定
	config := &genai.GenerateContentConfig{
		Temperature:     genai.Ptr(float32(0.7)),
		MaxOutputTokens: 500,
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				{Text: JudgeImagePrompt},
			},
		},
	}

	// チャットセッション作成
	chat, err := c.client.Chats.Create(ctx, c.defaultModel, config, nil)
	if err != nil {
		return JudgeResult{}, fmt.Errorf("failed to create chat: %w", err)
	}

	// 画像URLとテキストプロンプトを組み合わせて送信
	resp, err := chat.SendMessage(ctx,
		genai.Part{Text: "この写真を評価してください。"},
		genai.Part{FileData: &genai.FileData{
			MIMEType: "image/jpeg", // 一般的な画像形式として設定
			FileURI:  imageURL,
		}},
	)
	if err != nil {
		return JudgeResult{}, fmt.Errorf("failed to send message: %w", err)
	}

	// レスポンステキストを取得
	responseText := resp.Text()
	if responseText == "" {
		return JudgeResult{}, errors.New("empty response from AI")
	}

	// JSON解析を試みる
	var jsonResult struct {
		Score  int    `json:"score"`
		Reason string `json:"reason"`
	}

	// JSONとしてパース（簡易実装、エラーハンドリングは要改善）
	if err := parseJSONResponse(responseText, &jsonResult); err != nil {
		return JudgeResult{}, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return JudgeResult{
		Score:  jsonResult.Score,
		Reason: jsonResult.Reason,
	}, nil
}

// parseJSONResponseはAIのレスポンスからJSON部分を抽出してパースする
func parseJSONResponse(text string, v any) error {
	// AIが余分なテキストを含めることがあるため、JSON部分のみを抽出
	text = strings.TrimSpace(text)

	// コードブロック（```json ... ```）で囲まれている場合は除去
	if after, found := strings.CutPrefix(text, "```json"); found {
		text = strings.TrimSuffix(after, "```")
		text = strings.TrimSpace(text)
	} else if after, found := strings.CutPrefix(text, "```"); found {
		text = strings.TrimSuffix(after, "```")
		text = strings.TrimSpace(text)
	}

	// JSON部分をパース
	if err := json.Unmarshal([]byte(text), v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

// CloseはGeminiクライアントのクローズ処理（明示的なクローズは不要）
func (c *GeminiClient) Close() error {
	// 新しいgenaiクライアントは明示的なクローズ不要
	return nil
}
