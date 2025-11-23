package handlers

import (
	"context"
	"log"

	"github.com/ABfry/album-battler/backend/internal/domain/event"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
)

type RoomSettingsUpdatedHandler struct {
	eventPublisher service.EventPublisher
}

func NewRoomSettingsUpdatedHandler(
	eventPublisher service.EventPublisher,
) *RoomSettingsUpdatedHandler {
	return &RoomSettingsUpdatedHandler{
		eventPublisher: eventPublisher,
	}
}

// RoomSettingsUpdatedEventを処理
func (h *RoomSettingsUpdatedHandler) Handle(ctx context.Context, evt event.DomainEvent) error {
	e, ok := evt.(event.RoomSettingsUpdatedEvent)
	if !ok {
		log.Printf("Invalid event type: expected RoomSettingsUpdatedEvent, got %T", evt)
		return nil
	}

	// WebSocket通知を部屋のメンバーに送信
	broadcastEvent := service.BroadcastEvent{
		Type: e.EventType(),
		Payload: map[string]interface{}{
			"room_id":                   e.RoomID.String(),
			"room_number":               e.RoomNumber,
			"battle_time_limit_seconds": e.BattleTimeLimitSeconds,
		},
	}

	if err := h.eventPublisher.PublishToRoom(ctx, e.RoomID, broadcastEvent); err != nil {
		log.Printf("Failed to publish room_settings_updated event: %v", err)
		// WebSocket送信エラーは処理を止めない
	}

	log.Printf("Room settings updated in room %s (battle time limit: %ds)", e.RoomID, e.BattleTimeLimitSeconds)

	return nil
}
