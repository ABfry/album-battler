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

type UserInfo struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	IconURL string    `json:"icon_url"`
}

type GetRoomOutput struct {
	RoomNumber             int               `json:"room_number"`
	Users                  []UserInfo        `json:"users"`
	IsExpired              bool              `json:"is_expired"`
	RoomStatus             entity.RoomStatus `json:"room_status"`
	HostUserID             *uuid.UUID        `json:"host_user_id"`
	BattleTimeLimitSeconds int               `json:"battle_time_limit_seconds"`
}

type GetRoomUseCase struct {
	roomRepo repository.RoomRepository
	userRepo repository.UserRepository
}

func NewGetRoomUseCase(
	roomRepo repository.RoomRepository,
	userRepo repository.UserRepository,
) *GetRoomUseCase {
	return &GetRoomUseCase{
		roomRepo: roomRepo,
		userRepo: userRepo,
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

	// ユーザー情報を取得
	users, err := uc.userRepo.FindByIDs(ctx, room.UserIDs)
	if err != nil {
		return nil, err
	}

	// entity.User → UserInfoに変換
	userInfos := make([]UserInfo, len(users))
	for i, user := range users {
		userInfos[i] = UserInfo{
			ID:      user.ID,
			Name:    user.Name,
			IconURL: user.IconUrl,
		}
	}

	return &GetRoomOutput{
		RoomNumber:             room.RoomNumber,
		Users:                  userInfos,
		IsExpired:              room.IsExpired(),
		RoomStatus:             room.Status,
		HostUserID:             room.HostUserID,
		BattleTimeLimitSeconds: room.BattleTimeLimitSeconds,
	}, nil
}
