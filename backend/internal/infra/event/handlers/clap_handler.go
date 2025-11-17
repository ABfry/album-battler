package handlers

import (
	"context"
	"log"

	"github.com/ABfry/album-battler/backend/internal/domain/event"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
)

type ClapSendHandler struct {
	eventPublisher service.EventPublisher
}

func NewClapSendHandler(eventPublisher service.EventPublisher) *ClapSendHandler {
	return &ClapSendHandler{
		eventPublisher: eventPublisher,
	}
}

// ClapSendEventを処理
func (h *ClapSendHandler) Handle(ctx context.Context, evt event.DomainEvent) error {
	e, ok := evt.(event.ClapSendEvent)
	if !ok {
		log.Printf("Invalid event type: expected ClapSendEvent, got %T", evt)
		return nil
	}

	// WebSocket通知を部屋のメンバーに送信
	broadcastEvent := service.BroadcastEvent{
		Type:    e.EventType(),
		Payload: map[string]interface{}{"user_id": e.UserID, "clap_count": e.ClapCount},
	}

	if err := h.eventPublisher.PublishToRoom(ctx, e.RoomID, broadcastEvent); err != nil {
		log.Printf("Failed to publish clap_send event: %v", err)
		// WebSocket送信エラーは処理を止めない
	}

	log.Printf("Clap sent by user %s in battle %s", e.UserID, e.BattleID)
	return nil
}
