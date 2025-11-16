package handlers

import (
	"context"
	"log"

	"github.com/ABfry/album-battler/backend/internal/domain/event"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
)

type GameStartedHandler struct {
	eventPublisher service.EventPublisher
}

func NewGameStartedHandler(
	eventPublisher service.EventPublisher,
) *GameStartedHandler {
	return &GameStartedHandler{
		eventPublisher: eventPublisher,
	}
}

// GameStartedEventを処理
func (h *GameStartedHandler) Handle(ctx context.Context, evt event.DomainEvent) error {
	e, ok := evt.(event.GameStartedEvent)
	if !ok {
		log.Printf("Invalid event type: expected GameStartedEvent, got %T", evt)
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
