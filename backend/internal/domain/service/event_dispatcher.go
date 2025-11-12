package service

import (
	"context"

	"github.com/ABfry/album-battler/backend/internal/domain/event"
)

type EventDispatcher interface {
	// ドメインイベントを対応するハンドラーに配信
	Dispatch(ctx context.Context, events []event.DomainEvent) error
}
