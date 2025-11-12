# AI Client Infrastructure

このパッケージは、OpenAIとGemini APIを統一されたインターフェースで利用するための基盤を提供します。

## 概要

- **Interface定義**: `Client` インターフェースで両方のAIプロバイダーを抽象化
- **OpenAI実装**: `OpenAIClient` で OpenAI API をサポート
- **Gemini実装**: `GeminiClient` で Google Gemini API をサポート

## インストール

必要な依存関係は既に `go.mod` に含まれています：

```bash
cd backend
go mod download
```

## 使用方法

### 基本的なインターフェース

```go
type Client interface {
    Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)
    Close() error
}
```

### OpenAI クライアント

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/ABfry/album-battler/backend/internal/infrastructure/ai"
)

func main() {
    // クライアントの作成
    client, err := ai.NewOpenAIClient(ai.OpenAIConfig{
        APIKey:       "your-openai-api-key",
        DefaultModel: "gpt-4", // または "gpt-4o", "gpt-3.5-turbo"
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // テキスト生成
    resp, err := client.Generate(context.Background(), &ai.GenerateRequest{
        Messages: []ai.Message{
            {Role: "system", Content: "あなたは親切なアシスタントです。"},
            {Role: "user", Content: "Goについて教えてください。"},
        },
        MaxTokens:   1000,
        Temperature: 0.7,
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Response:", resp.Content)
    fmt.Printf("Tokens used: %d\n", resp.Usage.TotalTokens)
}
```

### Gemini クライアント

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/ABfry/album-battler/backend/internal/infrastructure/ai"
)

func main() {
    ctx := context.Background()

    // クライアントの作成
    client, err := ai.NewGeminiClient(ctx, ai.GeminiConfig{
        APIKey:       "your-gemini-api-key",
        DefaultModel: "gemini-2.0-flash", // または "gemini-1.5-pro"
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // テキスト生成
    resp, err := client.Generate(ctx, &ai.GenerateRequest{
        Messages: []ai.Message{
            {Role: "system", Content: "あなたは親切なアシスタントです。"},
            {Role: "user", Content: "Goについて教えてください。"},
        },
        MaxTokens:   1000,
        Temperature: 0.7,
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Response:", resp.Content)
    fmt.Printf("Tokens used: %d\n", resp.Usage.TotalTokens)
}
```

### プロバイダーの切り替え

同じインターフェースを使用しているため、簡単に切り替えができます：

```go
func getAIClient(ctx context.Context, provider string) (ai.Client, error) {
    switch provider {
    case "openai":
        return ai.NewOpenAIClient(ai.OpenAIConfig{
            APIKey: os.Getenv("OPENAI_API_KEY"),
        })
    case "gemini":
        return ai.NewGeminiClient(ctx, ai.GeminiConfig{
            APIKey: os.Getenv("GEMINI_API_KEY"),
        })
    default:
        return nil, fmt.Errorf("unknown provider: %s", provider)
    }
}
```

## データ構造

### GenerateRequest

```go
type GenerateRequest struct {
    Messages    []Message  // 会話履歴
    MaxTokens   int        // 最大トークン数（オプション）
    Temperature float32    // 温度パラメータ（オプション）
    Model       string     // モデル名（オプション、デフォルトはクライアント設定を使用）
}
```

### Message

```go
type Message struct {
    Role    string // "user", "assistant", "system"
    Content string // メッセージ内容
}
```

### GenerateResponse

```go
type GenerateResponse struct {
    Content      string     // 生成されたテキスト
    FinishReason string     // 終了理由
    Usage        TokenUsage // トークン使用状況
}
```

### TokenUsage

```go
type TokenUsage struct {
    PromptTokens     int // プロンプトトークン数
    CompletionTokens int // 生成トークン数
    TotalTokens      int // 合計トークン数
}
```

## 環境変数

APIキーは環境変数で管理することを推奨します：

```bash
export OPENAI_API_KEY="your-openai-api-key"
export GEMINI_API_KEY="your-gemini-api-key"
```

## 依存関係

- **OpenAI**: `github.com/openai/openai-go` v1.12.0
- **Gemini**: `google.golang.org/genai` v1.33.0

## 注意事項

- Geminiクライアントの作成時は`context.Context`が必要です
- 両方のクライアントとも、使用後は必ず`Close()`を呼び出してください
- APIキーは環境変数や設定ファイルで安全に管理してください
