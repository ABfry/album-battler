package service

import (
	"errors"
	"math/rand/v2"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
)

// RoomNumberGenerator は使用可能な部屋番号を生成するドメインサービス。
// why: 複数のRoom entityを見て判定する必要があり、単一のEntityの責務を超えるため。
type RoomNumberGenerator interface {
	Generate(activeRooms []*entity.Room) (int, error)
}

// RoomNumberGeneratorImpl は RoomNumberGenerator の実装。
type RoomNumberGeneratorImpl struct{}

// NewRoomNumberGenerator は RoomNumberGenerator の実装を生成する。
func NewRoomNumberGenerator() RoomNumberGenerator {
	return &RoomNumberGeneratorImpl{}
}

// Generate は使用されていない部屋番号をランダムに生成する。
// why: 1000-9999の範囲で、アクティブな部屋が使用していない番号からランダムに選ぶビジネスルール。
func (g *RoomNumberGeneratorImpl) Generate(activeRooms []*entity.Room) (int, error) {
	// 使用中の部屋番号を集める
	usedNumbers := make(map[int]bool)
	for _, room := range activeRooms {
		if room.IsActive() {
			usedNumbers[room.RoomNumber] = true
		}
	}

	// 使用可能な部屋番号のスライスを作成
	var availableNumbers []int
	for i := 1000; i <= 9999; i++ {
		if !usedNumbers[i] {
			availableNumbers = append(availableNumbers, i)
		}
	}

	// 使用可能な番号がない場合はエラー
	if len(availableNumbers) == 0 {
		return 0, errors.New("no available room number")
	}

	// ランダムに1つ選択
	randomIndex := rand.IntN(len(availableNumbers))
	return availableNumbers[randomIndex], nil
}
