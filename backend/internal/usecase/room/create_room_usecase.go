package room

import (
	"context"
	"errors"
	"strings"
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
	roomRepo            repository.RoomRepository
	dispatcher          service.EventDispatcher
	roomNumberGenerator service.RoomNumberGenerator
}

func NewCreateRoomUseCase(
	roomRepo repository.RoomRepository,
	dispatcher service.EventDispatcher,
	roomNumberGenerator service.RoomNumberGenerator,
) *CreateRoomUseCase {
	return &CreateRoomUseCase{
		roomRepo:            roomRepo,
		dispatcher:          dispatcher,
		roomNumberGenerator: roomNumberGenerator,
	}
}

func (uc *CreateRoomUseCase) Execute(ctx context.Context, input CreateRoomInput) (*CreateRoomOutput, error) {
	const maxRetries = 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		output, err := uc.tryCreateRoom(ctx, input)

		if err == nil {
			return output, nil // 成功
		}

		// 重複エラーの場合はリトライ
		if isDuplicateRoomNumberError(err) {
			continue
		}

		// それ以外のエラーは即座に返す
		return nil, err
	}

	return nil, errors.New("failed to create room after retries: no available room number")
}

func (uc *CreateRoomUseCase) tryCreateRoom(ctx context.Context, input CreateRoomInput) (*CreateRoomOutput, error) {
	// 既存の部屋を取得
	rooms, err := uc.roomRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// ドメインサービスで使用可能な部屋番号を生成
	roomNumber, err := uc.roomNumberGenerator.Generate(rooms)
	if err != nil {
		return nil, err
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

	// 永続化 (重複エラーが発生する可能性あり)
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

// 部屋番号の重複エラーかどうかを判定する
func isDuplicateRoomNumberError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	// MySQL duplicate entry error (1062)
	return strings.Contains(errStr, "Duplicate entry") ||
		strings.Contains(errStr, "duplicate key") ||
		strings.Contains(errStr, "1062")
}
