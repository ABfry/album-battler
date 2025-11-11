package room

import (
	"context"
	"errors"
	"time"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
	"github.com/google/uuid"
)

type CreateRoomInput struct {
	UserID uuid.UUID
}

type CreateRoomOutput struct {
	RoomID     uuid.UUID
	RoomNumber int
}

type CreateRoomUseCase struct {
	roomRepo   repository.RoomRepository
	dispatcher service.EventDispatcher
}

func NewCreateRoomUseCase(
	roomRepo repository.RoomRepository,
	dispatcher service.EventDispatcher,
) *CreateRoomUseCase {
	return &CreateRoomUseCase{
		roomRepo:   roomRepo,
		dispatcher: dispatcher,
	}
}

func (uc *CreateRoomUseCase) Execute(ctx context.Context, input CreateRoomInput) (*CreateRoomOutput, error) {
	// 既存の部屋を取得して、使用されていない部屋番号を見つける
	rooms, err := uc.roomRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// 使用中の部屋番号を集める
	usedNumbers := make(map[int]bool)
	for _, room := range rooms {
		usedNumbers[room.RoomNumber] = true
	}

	// 1~9999の範囲で使用されていない部屋番号を見つける
	var roomNumber int
	for i := 1; i <= 9999; i++ {
		if !usedNumbers[i] {
			roomNumber = i
			break
		}
	}

	if roomNumber == 0 {
		return nil, errors.New("no available room number")
	}

	// TODO 決める
	// 部屋の有効期限: 一旦24時間
	expiredAt := time.Now().Add(24 * time.Hour)

	// TODO 決めて 入力どう受け取るか決める
	// 部屋を作成 (maxUsers = 5)
	room, err := entity.NewRoom(roomNumber, expiredAt, 5)
	if err != nil {
		return nil, err
	}

	// 作成者を部屋に追加
	if err := room.AddUser(input.UserID); err != nil {
		return nil, err
	}

	// 永続化
	if err := uc.roomRepo.Save(ctx, room); err != nil {
		return nil, err
	}

	// ドメインイベントを配信
	events := room.PopEvents()
	if err := uc.dispatcher.Dispatch(ctx, events); err != nil {
		return nil, err
	}

	return &CreateRoomOutput{
		RoomID:     room.ID,
		RoomNumber: room.RoomNumber,
	}, nil
}
