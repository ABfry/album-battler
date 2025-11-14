package battle

import (
	"context"
	"errors"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/google/uuid"
)

// BattleIDを受け取り、Battleを返す
type GetBattleInput struct {
	BattleID uuid.UUID
}

type GetBattleOutput struct {
	Battle *entity.Battle
}

type GetBattleUseCase struct {
	battleRepo repository.BattleRepository
}

func NewGetBattleUseCase(battleRepo repository.BattleRepository) *GetBattleUseCase {
	return &GetBattleUseCase{battleRepo: battleRepo}
}

func (uc *GetBattleUseCase) Execute(ctx context.Context, input GetBattleInput) (*GetBattleOutput, error) {
	// バトルを取得
	battle, err := uc.battleRepo.FindByID(ctx, input.BattleID)
	if err != nil {
		return nil, errors.New("failed to get battle")
	}
	if battle == nil {
		return nil, errors.New("battle not found")
	}
	return &GetBattleOutput{Battle: battle}, nil
}
