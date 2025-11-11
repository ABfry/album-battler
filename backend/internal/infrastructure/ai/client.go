package ai

import "context"

// Messageは会話中の単一メッセージを表す
type Message struct {
	Role    string // "user", "assistant", "system" など
	Content string
}

// GenerateRequestはテキスト生成リクエストを表す
type GenerateRequest struct {
	Messages    []Message
	MaxTokens   int
	Temperature float32
	Model       string // 任意: 利用するモデルを指定
}

// GenerateResponseはテキスト生成のレスポンスを表す
type GenerateResponse struct {
	Content      string
	FinishReason string // "stop", "length", "content_filter" など
	Usage        TokenUsage
}

// TokenUsageはトークン使用情報を表す
type TokenUsage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// ClientはAIクライアント実装用のインターフェース
type Client interface {
	// Generateは与えられたメッセージに基づいてテキスト生成を行う
	Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)

	// Closeはクライアントをクローズしリソースを解放する
	Close() error
}
