package handlers

import (
	"context"
	"log"

	"github.com/ABfry/album-battler/backend/internal/domain/event"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
)

type StartButtonPressedHandler struct {
	eventPublisher service.EventPublisher
}

func NewStartButtonPressedHandler(
	eventPublisher service.EventPublisher,
) *StartButtonPressedHandler {
	return &StartButtonPressedHandler{
		eventPublisher: eventPublisher,
	}
}

// GameStartButtonPressedEventを処理
func (h *StartButtonPressedHandler) Handle(ctx context.Context, evt event.DomainEvent) error {
	e, ok := evt.(event.GameStartButtonPressedEvent)
	if !ok {
		log.Printf("Invalid event type: expected StartButtonPressedEvent, got %T", evt)
		return nil
	}

	// WebSocket通知を部屋のメンバーに送信
	broadcastEvent := service.BroadcastEvent{
		Type: e.EventType(),
		Payload: map[string]interface{}{
			"room_id": e.RoomID.String(),
		},
	}

	if err := h.eventPublisher.PublishToRoom(ctx, e.RoomID, broadcastEvent); err != nil {
		log.Printf("Failed to publish start_game event: %v", err)
		// WebSocket送信エラーは処理を止めない
	}

	log.Printf("Game started in room %s", e.RoomID)

	return nil
}
