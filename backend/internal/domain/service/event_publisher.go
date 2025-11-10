package service

import (
	"context"

	"github.com/google/uuid"
)

// 全ユーザーにブロードキャストするイベント
type BroadcastEvent struct {
	Type    string      // イベントタイプ
	Payload interface{} // ペイロード(データの中身)
}

// UseCase からイベントを発行するためのインターフェース
type EventPublisher interface {
	// 全接続ユーザーにイベントをブロードキャスト
	BroadcastToAll(ctx context.Context, event BroadcastEvent) error

	// 特定ユーザーにイベントを送信
	PublishToUser(ctx context.Context, userID uuid.UUID, event BroadcastEvent) error

	// 特定ルームにイベントを送信
	PublishToRoom(ctx context.Context, roomID uuid.UUID, event BroadcastEvent) error
}
