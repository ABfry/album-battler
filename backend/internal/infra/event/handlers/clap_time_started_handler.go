package handlers

import (
	"context"
	"log"

	"github.com/ABfry/album-battler/backend/internal/domain/event"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
)

type ClapTimeStartedHandler struct {
	eventPublisher service.EventPublisher
}

func NewClapTimeStartedHandler(eventPublisher service.EventPublisher) *ClapTimeStartedHandler {
	return &ClapTimeStartedHandler{
		eventPublisher: eventPublisher,
	}
}

// StartClapTimeEventを処理
func (h *ClapTimeStartedHandler) Handle(ctx context.Context, evt event.DomainEvent) error {
	e, ok := evt.(event.StartClapTimeEvent)
	if !ok {
		log.Printf("Invalid event type: expected StartClapTimeEvent, got %T", evt)
		return nil
	}

	// WebSocket通知を部屋のメンバーに送信
	broadcastEvent := service.BroadcastEvent{
		Type:    e.EventType(),
		Payload: map[string]interface{}{},
	}

	if err := h.eventPublisher.PublishToRoom(ctx, e.RoomID, broadcastEvent); err != nil {
		log.Printf("Failed to publish start_clap_time event: %v", err)
		// WebSocket送信エラーは処理を止めない
	}

	log.Printf("Clap time started for battle %s", e.BattleID)
	return nil
}
