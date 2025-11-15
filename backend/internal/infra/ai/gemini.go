package ai

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ABfry/album-battler/backend/internal/domain/service/llm"

	"google.golang.org/genai"
)

var _ llm.LLMClient = (*GeminiClient)(nil)

type GeminiClient struct {
	client       *genai.Client
	defaultModel string
}

// Geminiクライアントの設定
type GeminiConfig struct {
	APIKey       string
	DefaultModel string // 例: "gemini-2.0-flash", "gemini-1.5-pro"
}

func NewGeminiClient(ctx context.Context, config GeminiConfig) (*GeminiClient, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("gemini API key is required")
	}

	if config.DefaultModel == "" {
		config.DefaultModel = "gemini-2.5-flash"
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

// GeminiのAPIでテキスト生成する
func (c *GeminiClient) Generate(ctx context.Context, req *llm.GenerateRequest) (*llm.GenerateResponse, error) {
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

	// MaxOutputTokensを設定
	// 注: Gemini SDKではMaxOutputTokensを設定しない場合、モデルのデフォルト値が使用される
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
	usage := llm.TokenUsage{}
	if resp.UsageMetadata != nil {
		usage.PromptTokens = int(resp.UsageMetadata.PromptTokenCount)
		usage.CompletionTokens = int(resp.UsageMetadata.CandidatesTokenCount)
		usage.TotalTokens = int(resp.UsageMetadata.TotalTokenCount)
	}

	return &llm.GenerateResponse{
		Content:      content,
		FinishReason: finishReason,
		Usage:        usage,
	}, nil
}

// 写真撮影バトル用のテーマを生成
func (c *GeminiClient) GenerateTheme(ctx context.Context) (string, error) {
	// Function Declarationを定義
	generateThemeFunc := &genai.FunctionDeclaration{
		Name:        "generate_theme",
		Description: "写真撮影バトル用のテーマを生成する",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"theme": {
					Type:        genai.TypeString,
					Description: "生成されたテーマ",
				},
			},
			Required: []string{"theme"},
		},
	}

	// プロンプトを構築
	parts := []*genai.Part{
		genai.NewPartFromText(GenerateThemePrompt),
	}

	// Contentsを構築
	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	// 生成設定（Function Calling有効化）
	config := &genai.GenerateContentConfig{
		Temperature: genai.Ptr(float32(0.9)),
		Tools: []*genai.Tool{
			{
				FunctionDeclarations: []*genai.FunctionDeclaration{generateThemeFunc},
			},
		},
	}

	// Models APIでテーマ生成を実行
	result, err := c.client.Models.GenerateContent(ctx, c.defaultModel, contents, config)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	// Function Callのレスポンスを処理
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("no response from AI")
	}

	// Function Callのパートを探す
	var functionCall *genai.FunctionCall
	for _, part := range result.Candidates[0].Content.Parts {
		if part.FunctionCall != nil {
			functionCall = part.FunctionCall
			break
		}
	}

	if functionCall == nil {
		return "", errors.New("no function call in response")
	}

	// Function Callの引数からテーマを取得
	themeArg, ok := functionCall.Args["theme"]
	if !ok {
		return "", errors.New("theme not found in function call")
	}

	// 型変換
	theme, ok := themeArg.(string)
	if !ok {
		return "", fmt.Errorf("invalid theme type: %T", themeArg)
	}

	if theme == "" {
		return "", errors.New("generated theme is empty")
	}

	return theme, nil
}

// 提出された画像を評価し、スコアと理由を返す
func (c *GeminiClient) JudgeImage(ctx context.Context, req *llm.ImageJudgeRequest) (*llm.ImageJudgeResponse, error) {
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

// downloadImageはURLから画像をダウンロードしてバイト配列とMIMEタイプを返す
func downloadImage(ctx context.Context, url string) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download image: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close response body: %v\n", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to download image: status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read image data: %w", err)
	}

	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "image/jpeg" // デフォルト
	}

	return data, mimeType, nil
}

// 1枚の画像を評価する内部ヘルパー関数（Models API使用）
func (c *GeminiClient) judgeOneImage(ctx context.Context, theme, imageURL string) (llm.JudgeResult, error) {
	// 画像をダウンロード
	imageBytes, mimeType, err := downloadImage(ctx, imageURL)
	if err != nil {
		return llm.JudgeResult{}, fmt.Errorf("failed to download image: %w", err)
	}

	// Function Declarationを定義
	judgeImageFunc := &genai.FunctionDeclaration{
		Name:        "judge_image",
		Description: "画像を評価してスコア（0-100）と理由を返す",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"score": {
					Type:        genai.TypeInteger,
					Description: "画像の評価スコア（0-100）",
				},
				"reason": {
					Type:        genai.TypeString,
					Description: "評価の理由",
				},
			},
			Required: []string{"score", "reason"},
		},
	}

	// プロンプトテキストと画像のPartsを構築
	prompt := fmt.Sprintf("お題「%s」に対して、この写真を評価してください。", theme)
	parts := []*genai.Part{
		genai.NewPartFromText(prompt),
		genai.NewPartFromBytes(imageBytes, mimeType),
	}

	// Contentsを構築
	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	// 生成設定（Function Calling有効化）
	config := &genai.GenerateContentConfig{
		Temperature: genai.Ptr(float32(0.7)),
		// MaxOutputTokensを設定しない（APIデフォルト値を使用）
		SystemInstruction: genai.NewContentFromParts([]*genai.Part{
			genai.NewPartFromText(JudgeImagePrompt),
		}, genai.RoleUser),
		Tools: []*genai.Tool{
			{
				FunctionDeclarations: []*genai.FunctionDeclaration{judgeImageFunc},
			},
		},
	}

	// Models APIで画像評価を実行
	result, err := c.client.Models.GenerateContent(ctx, c.defaultModel, contents, config)
	if err != nil {
		return llm.JudgeResult{}, fmt.Errorf("failed to generate content: %w", err)
	}

	// Function Callのレスポンスを処理
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return llm.JudgeResult{}, errors.New("no response from AI")
	}

	// Function Callのパートを探す
	var functionCall *genai.FunctionCall
	for _, part := range result.Candidates[0].Content.Parts {
		if part.FunctionCall != nil {
			functionCall = part.FunctionCall
			break
		}
	}

	if functionCall == nil {
		return llm.JudgeResult{}, errors.New("no function call in response")
	}

	// Function Callの引数からスコアと理由を取得
	scoreArg, ok := functionCall.Args["score"]
	if !ok {
		return llm.JudgeResult{}, errors.New("score not found in function call")
	}

	reasonArg, ok := functionCall.Args["reason"]
	if !ok {
		return llm.JudgeResult{}, errors.New("reason not found in function call")
	}

	// 型変換
	score, ok := scoreArg.(float64)
	if !ok {
		return llm.JudgeResult{}, fmt.Errorf("invalid score type: %T", scoreArg)
	}

	reason, ok := reasonArg.(string)
	if !ok {
		return llm.JudgeResult{}, fmt.Errorf("invalid reason type: %T", reasonArg)
	}

	return llm.JudgeResult{
		Score:  int(score),
		Reason: reason,
	}, nil
}

// Geminiクライアントのクローズ処理（明示的なクローズは不要）
func (c *GeminiClient) Close() error {
	// 明示的なクローズ不要
	return nil
}
