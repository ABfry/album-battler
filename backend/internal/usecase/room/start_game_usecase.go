package room

import (
	"context"
	"errors"

	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
	"github.com/google/uuid"
)

type StartGameInput struct {
	RoomID uuid.UUID
	UserID uuid.UUID // ホスト検証用
}

type StartGameUseCase struct {
	roomRepo   repository.RoomRepository
	dispatcher service.EventDispatcher
}

func NewStartGameUseCase(
	roomRepo repository.RoomRepository,
	dispatcher service.EventDispatcher,
) *StartGameUseCase {
	return &StartGameUseCase{
		roomRepo:   roomRepo,
		dispatcher: dispatcher,
	}
}

func (uc *StartGameUseCase) Execute(ctx context.Context, input StartGameInput) error {
	// 部屋を取得
	room, err := uc.roomRepo.FindByID(ctx, input.RoomID)
	if err != nil {
		return err
	}

	// ホストであることを確認
	if room.HostUserID == nil || *room.HostUserID != input.UserID {
		return errors.New("only host can start the game")
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
