package service

import (
	"context"
	"io"
)

// 画像ファイルのストレージ操作を抽象化するインターフェース
type ImageStorage interface {
	// 画像をアップロードし、アクセス可能なURLを返す
	// key: ストレージ内での一意なキー
	// data: アップロードする画像データ
	Upload(ctx context.Context, key string, data io.Reader, contentType string) (string, error)

	// 指定されたキーの画像のアクセス用URLを取得する
	GetURL(ctx context.Context, key string) (string, error)
}
