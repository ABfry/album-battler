package websocket

import (
	"github.com/google/uuid"

	"github.com/ABfry/album-battler/backend/internal/domain/service"
)

type WebSocketRoomManager struct {
	hub *Hub
}

func NewWebSocketRoomManager(hub *Hub) service.RoomManager {
	return &WebSocketRoomManager{hub: hub}
}

func (rm *WebSocketRoomManager) JoinRoom(userID, roomID uuid.UUID) {
	rm.hub.JoinRoom(userID, roomID)
}

func (rm *WebSocketRoomManager) LeaveRoom(userID, roomID uuid.UUID) {
	rm.hub.LeaveRoom(userID, roomID)
}

func (rm *WebSocketRoomManager) GetRoomMembers(roomID uuid.UUID) []uuid.UUID {
	return rm.hub.GetRoomMembers(roomID)
}
