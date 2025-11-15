package battle

import (
	"context"
	"errors"
	"fmt"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type CreateBattleInput struct {
	RoomID uuid.UUID
}

type CreateBattleOutput struct {
	BattleID uuid.UUID
}

type CreateBattleUseCase struct {
	battleRepo     repository.BattleRepository
	battleUserRepo repository.BattleUserRepository
	roomRepo       repository.RoomRepository
}

func NewCreateBattleUseCase(
	battleRepo repository.BattleRepository,
	battleUserRepo repository.BattleUserRepository,
	roomRepo repository.RoomRepository,
) *CreateBattleUseCase {
	return &CreateBattleUseCase{
		battleRepo:     battleRepo,
		battleUserRepo: battleUserRepo,
		roomRepo:       roomRepo,
	}
}

func (uc *CreateBattleUseCase) Execute(ctx context.Context, input CreateBattleInput) (*CreateBattleOutput, error) {

	// ルームから参加者UserID一覧を取得
	room, err := uc.roomRepo.FindByID(ctx, input.RoomID)
	if err != nil {
		return nil, errors.New("failed to find room")
	}
	userIDs := room.UserIDs

	// バトルを作成
	// todo: テーマAIから生成
	battle, err := entity.NewBattle(input.RoomID, userIDs, "test仮")
	if err != nil {
		return nil, errors.New("failed to create battle")
	}

	err = uc.battleRepo.Save(ctx, battle)
	if err != nil {
		return nil, errors.New("failed to save battle")
	}

	// battle_users 中間テーブルに参加者を一括保存
	if err := uc.battleUserRepo.SaveBatch(ctx, battle.ID, userIDs); err != nil {
		fmt.Println("failed to save battle users", err)
		return nil, errors.New("failed to save battle users")
	}

	return &CreateBattleOutput{BattleID: battle.ID}, nil
}
