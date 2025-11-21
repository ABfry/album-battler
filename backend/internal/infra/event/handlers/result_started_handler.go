package handlers

import (
	"context"
	"log"

	"github.com/ABfry/album-battler/backend/internal/domain/event"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
)

type ResultStartedHandler struct {
	eventPublisher service.EventPublisher
}

func NewResultStartedHandler(eventPublisher service.EventPublisher) *ResultStartedHandler {
	return &ResultStartedHandler{
		eventPublisher: eventPublisher,
	}
}

// StartResultPhaseEventを処理
func (h *ResultStartedHandler) Handle(ctx context.Context, evt event.DomainEvent) error {
	e, ok := evt.(event.StartResultPhaseEvent)
	if !ok {
		log.Printf("Invalid event type: expected StartResultPhaseEvent, got %T", evt)
		return nil
	}

	// WebSocket通知を部屋のメンバーに送信（イベント発生の通知のみ）
	broadcastEvent := service.BroadcastEvent{
		Type:    e.EventType(),
		Payload: map[string]interface{}{},
	}

	if err := h.eventPublisher.PublishToRoom(ctx, e.RoomID, broadcastEvent); err != nil {
		log.Printf("Failed to publish start_result_phase event: %v", err)
		// WebSocket送信エラーは処理を止めない
	}

	log.Printf("Result phase started for battle %s, winner: %s", e.BattleID, e.WinnerUserID)
	return nil
}
