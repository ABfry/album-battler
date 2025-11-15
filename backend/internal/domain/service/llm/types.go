package llm

// Message は会話中の単一メッセージを表す
type Message struct {
	Role    string // "user", "assistant", "system" など
	Content string
}

// GenerateRequest はテキスト生成リクエストを表す
type GenerateRequest struct {
	Messages    []Message
	MaxTokens   int
	Temperature float32
	Model       string // 任意: 利用するモデルを指定
}

// GenerateResponse はテキスト生成のレスポンスを表す
type GenerateResponse struct {
	Content      string
	FinishReason string // "stop", "length", "content_filter" など
	Usage        TokenUsage
}

// TokenUsage はトークン使用情報を表す
type TokenUsage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// ImageJudgeRequest は画像評価リクエストを表す
type ImageJudgeRequest struct {
	Theme     string   // お題
	ImageURLs []string // 評価する画像のURL一覧
}

// ImageJudgeResponse は画像評価レスポンスを表す
type ImageJudgeResponse struct {
	Results []JudgeResult
}

// JudgeResult は個別の画像評価結果を表す
type JudgeResult struct {
	Score  int    // 評価スコア（0-100）
	Reason string // 評価理由
}
