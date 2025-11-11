package event

import (
	"context"
	"fmt"
	"log"

	"github.com/ABfry/album-battler/backend/internal/domain/event"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
)

type EventHandler interface {
	Handle(ctx context.Context, event event.DomainEvent) error
}

type EventDispatcherImpl struct {
	handlers map[string][]EventHandler
}

func NewEventDispatcher() service.EventDispatcher {
	return &EventDispatcherImpl{
		handlers: make(map[string][]EventHandler),
	}
}

// 特定のイベントタイプに対するハンドラーを登録
func (d *EventDispatcherImpl) Register(eventType string, handler EventHandler) {
	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

// イベントを対応するハンドラーに配信
func (d *EventDispatcherImpl) Dispatch(ctx context.Context, events []event.DomainEvent) error {
	for _, evt := range events {
		eventType := evt.EventType()
		handlers, ok := d.handlers[eventType]
		if !ok {
			log.Printf("No handler registered for event type: %s", eventType)
			continue
		}

		for _, handler := range handlers {
			if err := handler.Handle(ctx, evt); err != nil {
				log.Printf("Error handling event %s: %v", eventType, err)
				return fmt.Errorf("failed to handle event %s: %w", eventType, err)
			}
		}
	}

	return nil
}
