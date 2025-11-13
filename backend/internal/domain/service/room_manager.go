package service

import "github.com/google/uuid"

type RoomManager interface {
	JoinRoom(userID uuid.UUID, roomID uuid.UUID)
	LeaveRoom(userID uuid.UUID, roomID uuid.UUID)
	GetRoomMembers(roomID uuid.UUID) []uuid.UUID
}
