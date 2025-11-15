package room

import (
	"context"
	"errors"
	"fmt"

	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
	battleusecase "github.com/ABfry/album-battler/backend/internal/usecase/battle"
	"github.com/google/uuid"
)

type StartGameInput struct {
	RoomID uuid.UUID
	UserID uuid.UUID // ホスト検証用
}

type StartGameUseCase struct {
	roomRepo      repository.RoomRepository
	dispatcher    service.EventDispatcher
	battleUseCase BattleCreator
}

type BattleCreator interface {
	Execute(ctx context.Context, input battleusecase.CreateBattleInput) (*battleusecase.CreateBattleOutput, error)
}

func NewStartGameUseCase(
	roomRepo repository.RoomRepository,
	dispatcher service.EventDispatcher,
	battleUseCase BattleCreator,
) *StartGameUseCase {
	return &StartGameUseCase{
		roomRepo:      roomRepo,
		dispatcher:    dispatcher,
		battleUseCase: battleUseCase,
	}
}

func (uc *StartGameUseCase) Execute(ctx context.Context, input StartGameInput) error {
	// 部屋を取得
	room, err := uc.roomRepo.FindByID(ctx, input.RoomID)
	if err != nil {
		return err
	}
	if room == nil {
		return errors.New("room not found")
	}

	// ホストであることを確認
	if room.HostUserID == nil || *room.HostUserID != input.UserID {
		return errors.New("only host can start the game")
	}

	// ゲーム開始と同時にバトルを生成
	if _, err := uc.battleUseCase.Execute(ctx, battleusecase.CreateBattleInput{
		RoomID: room.ID,
	}); err != nil {
		return fmt.Errorf("failed to create battle: %w", err)
	}

	// ゲーム開始 (状態遷移 + イベント記録)
	if err := room.StartGame(); err != nil {
		return err
	}

	// 永続化
	if err := uc.roomRepo.Save(ctx, room); err != nil {
		return err
	}

	// ドメインイベントを配信
	events := room.PopEvents()
	if err := uc.dispatcher.Dispatch(ctx, events); err != nil {
		return err
	}

	return nil
}
