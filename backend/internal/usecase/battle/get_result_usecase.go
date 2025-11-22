package battle

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type GetResultInput struct {
	BattleID uuid.UUID
}

type UserResult struct {
	UserID        uuid.UUID `json:"user_id"`
	AIScore       float64   `json:"ai_score"`
	UserScore     int       `json:"user_score"`
	FinalScore    float64   `json:"final_score"`
	Rank          int       `json:"rank"`
	ImageURL      string    `json:"image_url"`
	AIExplanation string    `json:"ai_explanation"`
}

type GetResultOutput struct {
	BattleID     uuid.UUID    `json:"battle_id"`
	WinnerUserID uuid.UUID    `json:"winner_user_id"`
	Results      []UserResult `json:"results"`
}

type GetResultUseCase struct {
	battleRepo repository.BattleRepository
	roomRepo   repository.RoomRepository
	imageRepo  repository.ImageRepository
}

func NewGetResultUseCase(
	battleRepo repository.BattleRepository,
	roomRepo repository.RoomRepository,
	imageRepo repository.ImageRepository,
) *GetResultUseCase {
	return &GetResultUseCase{
		battleRepo: battleRepo,
		roomRepo:   roomRepo,
		imageRepo:  imageRepo,
	}
}

func (uc *GetResultUseCase) Execute(ctx context.Context, input GetResultInput) (*GetResultOutput, error) {
	fmt.Printf("GetResultUseCase Execute: battleID=%s\n", input.BattleID)
	// バトルを取得
	battle, err := uc.battleRepo.FindByID(ctx, input.BattleID)
	if err != nil {
		fmt.Printf("Failed to find battle: %v", err)
		return nil, errors.New("failed to find battle")
	}
	if battle == nil {
		fmt.Printf("Battle not found: %s", input.BattleID)
		return nil, errors.New("battle not found")
	}

	// Roomを取得してステータスを確認
	room, err := uc.roomRepo.FindByID(ctx, battle.RoomID)
	if err != nil {
		fmt.Printf("Failed to find room: %v", err)
		return nil, errors.New("failed to find room")
	}
	if room == nil {
		fmt.Printf("Room not found: %s", battle.RoomID)
		return nil, errors.New("room not found")
	}

	// Roomが結果フェーズでない場合はエラー
	if room.Status != entity.Result {
		fmt.Printf("Room is not in result phase: status=%v", room.Status)
		return nil, errors.New("results are not yet available")
	}

	// 全画像を取得
	images, err := uc.imageRepo.FindImagesByBattleID(ctx, input.BattleID)
	if err != nil {
		fmt.Printf("Failed to find images: %v", err)
		return nil, errors.New("failed to find images")
	}

	// スコア計算とソート用の構造体
	type scoreData struct {
		UserID        uuid.UUID
		AIScore       float64
		UserScore     int
		FinalScore    float64
		UploadedAt    time.Time
		ImageURL      string
		AIExplanation string
	}

	scores := make([]scoreData, 0, len(images))
	ranks := make(map[uuid.UUID]int)
	var winnerUserID uuid.UUID

	for _, img := range images {
		finalScore := img.AIScore + float64(img.UserScore)
		fmt.Printf("Image %s: AIScore=%.2f, UserScore=%d, FinalScore=%.2f\n", img.ID, img.AIScore, img.UserScore, finalScore)
		scores = append(scores, scoreData{
			UserID:        img.UserID,
			AIScore:       img.AIScore,
			UserScore:     img.UserScore,
			FinalScore:    finalScore,
			UploadedAt:    img.UploadedAt,
			ImageURL:      img.ImageURL,
			AIExplanation: img.AIExplanation,
		})
	}

	// 勝者決定＆順位付け
	sorted := make([]scoreData, len(scores))
	copy(sorted, scores)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].FinalScore != sorted[j].FinalScore {
			return sorted[i].FinalScore > sorted[j].FinalScore
		}

		return sorted[i].UploadedAt.Before(sorted[j].UploadedAt)
	})

	for i, score := range sorted {
		ranks[score.UserID] = i + 1
	}

	if len(sorted) > 0 {
		winnerUserID = sorted[0].UserID
	}

	// 結果をランク付きで返却
	results := make([]UserResult, 0, len(scores))
	for _, score := range scores {
		results = append(results, UserResult{
			UserID:        score.UserID,
			AIScore:       score.AIScore,
			UserScore:     score.UserScore,
			FinalScore:    score.FinalScore,
			Rank:          ranks[score.UserID],
			ImageURL:      score.ImageURL,
			AIExplanation: score.AIExplanation,
		})
	}

	return &GetResultOutput{
		BattleID:     input.BattleID,
		WinnerUserID: winnerUserID,
		Results:      results,
	}, nil
}
