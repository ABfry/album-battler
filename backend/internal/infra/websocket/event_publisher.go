package websocket

import (
	"context"
	"time"

	"github.com/ABfry/album-battler/backend/internal/domain/service"
	"github.com/google/uuid"
)

type WebSocketEventPublisher struct {
	hub *Hub
}

func NewWebSocketEventPublisher(hub *Hub) service.EventPublisher {
	return &WebSocketEventPublisher{
		hub: hub,
	}
}

func (p *WebSocketEventPublisher) BroadcastToAll(
	ctx context.Context,
	event service.BroadcastEvent,
) error {
	message := Message{
		Type:      event.Type,
		Payload:   event.Payload,
		Timestamp: time.Now(),
	}

	return p.hub.BroadcastToAll(message)
}

func (p *WebSocketEventPublisher) PublishToUser(
	ctx context.Context,
	userID uuid.UUID,
	event service.BroadcastEvent,
) error {
	message := Message{
		Type:      event.Type,
		Payload:   event.Payload,
		Timestamp: time.Now(),
	}

	return p.hub.SendToUser(userID, message)
}

func (p *WebSocketEventPublisher) PublishToRoom(
	ctx context.Context,
	roomID uuid.UUID,
	event service.BroadcastEvent,
) error {
	message := Message{
		Type:      event.Type,
		Payload:   event.Payload,
		Timestamp: time.Now(),
	}

	return p.hub.BroadcastToRoom(roomID, message)
}
