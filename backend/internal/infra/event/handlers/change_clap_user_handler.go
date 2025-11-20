package handlers

import (
	"context"
	"log"

	"github.com/ABfry/album-battler/backend/internal/domain/event"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
)

type ClapChangeUserHandler struct {
	eventPublisher service.EventPublisher
}

func NewClapChangeUserHandler(
	eventPublisher service.EventPublisher,
) *ClapChangeUserHandler {
	return &ClapChangeUserHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *ClapChangeUserHandler) Handle(ctx context.Context, evt event.DomainEvent) error {
	e, ok := evt.(event.ChangeClapUserEvent)
	if !ok {
		log.Printf("Invalid event type: expected ChangeClapUserEvent, got %T", evt)
		return nil
	}

	// WebSocket通知を部屋のメンバーに送信
	broadcastEvent := service.BroadcastEvent{
		Type:    e.EventType(),
		Payload: map[string]interface{}{},
	}

	if err := h.eventPublisher.PublishToRoom(ctx, e.RoomID, broadcastEvent); err != nil {
		log.Printf("Failed to publish change_clap_user event: %v", err)
		// WebSocket送信エラーは処理を止めない
	}

	log.Printf("Clap user changed to %s in battle %s", e.UserID, e.BattleID)
	return nil
}
