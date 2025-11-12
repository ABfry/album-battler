package room

import (
	"context"
	"errors"
	"fmt"

	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
	"github.com/google/uuid"
)

type JoinRoomInput struct {
	UserID     uuid.UUID
	RoomNumber int
}

type JoinRoomUseCase struct {
	roomRepo   repository.RoomRepository
	dispatcher service.EventDispatcher
}

func NewJoinRoomUseCase(
	roomRepo repository.RoomRepository,
	dispatcher service.EventDispatcher,
) *JoinRoomUseCase {
	return &JoinRoomUseCase{
		roomRepo:   roomRepo,
		dispatcher: dispatcher,
	}
}

func (uc *JoinRoomUseCase) Execute(ctx context.Context, input JoinRoomInput) error {
	// 部屋を取得
	room, err := uc.roomRepo.FindByRoomNumber(ctx, input.RoomNumber)
	if err != nil {
		return err
	}
	if room == nil {
		return errors.New("room not found")
	}

	// 期限チェック
	if room.IsExpired() {
		return errors.New("room has expired")
	}

	// 満員チェック
	if room.IsFull() {
		return errors.New("room is full")
	}

	// ユーザーを追加 (ドメインロジック + イベント記録)
	if err := room.AddUser(input.UserID); err != nil {
		return err
	}

	// 永続化
	if err := uc.roomRepo.Save(ctx, room); err != nil {
		return err
	}

	// ドメインイベントを配信
	events := room.PopEvents()
	if err := uc.dispatcher.Dispatch(ctx, events); err != nil {
		fmt.Println("failed to dispatch domain events", "error", err)
		return err
	}

	return nil
}
