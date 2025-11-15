package battle

import (
	"context"
	"errors"

	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/google/uuid"
)

// RoomIDを受け取り、BattleIDを返す
type GetBattleIDInput struct {
	RoomID uuid.UUID
}

type GetBattleIDOutput struct {
	BattleID uuid.UUID
}

type GetBattleIDUseCase struct {
	battleRepo repository.BattleRepository
}

func NewGetBattleIDUseCase(battleRepo repository.BattleRepository) *GetBattleIDUseCase {
	return &GetBattleIDUseCase{battleRepo: battleRepo}
}

func (uc *GetBattleIDUseCase) Execute(ctx context.Context, input GetBattleIDInput) (*GetBattleIDOutput, error) {
	// 部屋に紐づくバトルを取得
	battle, err := uc.battleRepo.FindByRoomID(ctx, input.RoomID)

	// 部屋が存在しない場合
	if err != nil {
		return nil, errors.New("failed to get battle")
	}

	// バトルが存在しない場合
	if battle == nil {
		return nil, errors.New("battle not found")
	}

	return &GetBattleIDOutput{BattleID: battle.ID}, nil
}
