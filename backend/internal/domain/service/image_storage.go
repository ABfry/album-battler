package service

import (
	"context"
	"io"
)

// ImageStorage は画像ファイルのストレージ操作を抽象化するインターフェース
// S3やCloudinaryなど、具体的なストレージ実装に依存しない
type ImageStorage interface {
	// Upload は画像をアップロードし、アクセス可能なURLを返す
	// key: ストレージ内での一意なキー（例: "albums/{albumID}/image.jpg"）
	// data: アップロードする画像データ
	Upload(ctx context.Context, key string, data io.Reader, contentType string) (string, error)

	// Delete は指定されたキーの画像を削除する
	Delete(ctx context.Context, key string) error

	// GetURL は指定されたキーの画像のアクセス用URLを取得する
	// 署名付きURLが必要な場合は、ここで生成される
	GetURL(ctx context.Context, key string) (string, error)
}
