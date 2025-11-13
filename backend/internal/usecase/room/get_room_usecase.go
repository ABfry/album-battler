package room

import (
	"context"
	"errors"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type GetRoomInput struct {
	RoomID uuid.UUID
}

type GetRoomOutput struct {
	RoomNumber int           `json:"room_number"`
	Users      []uuid.UUID   `json:"users"`
	IsExpired  bool          `json:"is_expired"`
	RoomStatus entity.RoomStatus `json:"room_status"`
	HostUserID *uuid.UUID    `json:"host_user_id"`
}

type GetRoomUseCase struct {
	roomRepo repository.RoomRepository
}

func NewGetRoomUseCase(
	roomRepo repository.RoomRepository,
) *GetRoomUseCase {
	return &GetRoomUseCase{
		roomRepo: roomRepo,
	}
}

func (uc *GetRoomUseCase) Execute(ctx context.Context, input GetRoomInput) (*GetRoomOutput, error) {
	room, err := uc.roomRepo.FindByID(ctx, input.RoomID)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, errors.New("room not found")
	}

	return &GetRoomOutput{
		RoomNumber: room.RoomNumber,
		Users:      room.UserIDs,
		IsExpired:  room.IsExpired(),
		RoomStatus: room.Status,
		HostUserID: room.HostUserID,
	}, nil
}
