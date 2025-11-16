package battle

import (
	"context"
	"errors"
	"fmt"

	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type GetImageInput struct {
	BattleID uuid.UUID
}

type ImageInfo struct {
	UserID   uuid.UUID
	ImageURL string
}

type GetImageOutput struct {
	Images []ImageInfo
}

type GetImageUseCase struct {
	battleRepo repository.BattleRepository
	imageRepo  repository.ImageRepository
}

func NewGetImageUseCase(battleRepo repository.BattleRepository, imageRepo repository.ImageRepository) *GetImageUseCase {
	return &GetImageUseCase{
		battleRepo: battleRepo,
		imageRepo:  imageRepo,
	}
}

func (uc *GetImageUseCase) Execute(ctx context.Context, input *GetImageInput) (*GetImageOutput, error) {
	// バトルを取得
	battle, err := uc.battleRepo.FindByID(ctx, input.BattleID)
	if err != nil {
		fmt.Printf("Failed to get battle %s: %v", input.BattleID.String(), err)
		return nil, errors.New("failed to get battle")
	}
	if battle == nil {
		return nil, errors.New("battle not found")
	}

	images, err := uc.imageRepo.FindImagesByBattleID(ctx, battle.ID)
	if err != nil {
		fmt.Printf("Failed to get images for battle %s: %v", battle.ID.String(), err)
		return nil, errors.New("failed to get images")
	}

	imageInfos := make([]ImageInfo, 0, len(images))
	for _, image := range images {
		imageInfos = append(imageInfos, ImageInfo{
			UserID:   image.UserID,
			ImageURL: image.ImageURL,
		})
	}

	return &GetImageOutput{Images: imageInfos}, nil
}
