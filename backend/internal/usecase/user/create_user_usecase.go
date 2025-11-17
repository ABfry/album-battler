package user

import (
	"context"
	"fmt"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type CreateUserInput struct {
	UserName string
}

type CreateUserOutput struct {
	UserID uuid.UUID
}

type CreateUserUseCase struct {
	userRepo repository.UserRepository
}

func NewCreateUserUseCase(userRepo repository.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{userRepo: userRepo}
}

// ユーザーを作成する
func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	// iconURLは仮
	iconURL := "https://album-battler-images.s3.ap-northeast-1.amazonaws.com/icon/default/icon.png"

	user, err := entity.NewUser(input.UserName, iconURL, "password") // todo: passwordは仮
	if err != nil {
		return nil, fmt.Errorf("failed to create user")
	}

	if err := uc.userRepo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user")
	}

	return &CreateUserOutput{
		UserID: user.ID,
	}, nil
}
