package llm

import "context"

// LLMClientはLLM（大規模言語モデル）サービスのインターフェース
type LLMClient interface {
	// Generateは与えられたメッセージに基づいてテキスト生成を行う
	Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)

	// GenerateThemeはBattleのお題を生成する
	GenerateTheme(ctx context.Context) (string, error)

	// JudgeImageは画像を受け取り、点数をジャッジする
	JudgeImage(ctx context.Context, req *ImageJudgeRequest) (*ImageJudgeResponse, error)

	// Closeはクライアントをクローズしリソースを解放する
	Close() error
}
