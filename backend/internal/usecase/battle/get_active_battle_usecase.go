package battle

import (
	"context"
	"errors"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type GetActiveBattleInput struct {
	UserID uuid.UUID
}

type GetActiveBattleOutput struct {
	RoomID   uuid.UUID
	BattleID uuid.UUID
}

type GetActiveBattleUseCase struct {
	roomRepo   repository.RoomRepository
	battleRepo repository.BattleRepository
}

func NewGetActiveBattleUseCase(
	roomRepo repository.RoomRepository,
	battleRepo repository.BattleRepository,
) *GetActiveBattleUseCase {
	return &GetActiveBattleUseCase{
		roomRepo:   roomRepo,
		battleRepo: battleRepo,
	}
}

func (uc *GetActiveBattleUseCase) Execute(ctx context.Context, input GetActiveBattleInput) (*GetActiveBattleOutput, error) {
	// ユーザーが参加しているルームを取得
	rooms, err := uc.roomRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, errors.New("failed to find rooms by user ID")
	}

	for _, room := range rooms {
		// クローズドステータスのルームはスキップ
		if room.Status == entity.Closed {
			continue
		}

		// バトルを取得
		battle, err := uc.battleRepo.FindByRoomID(ctx, room.ID)
		if err != nil {
			return nil, errors.New("failed to find battle by room ID")
		}
		if battle == nil {
			continue
		}

		return &GetActiveBattleOutput{RoomID: room.ID, BattleID: battle.ID}, nil
	}

	return nil, nil // アクティブなルームが見つからなかった場合
}
