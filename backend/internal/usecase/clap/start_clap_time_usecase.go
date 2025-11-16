package clap

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type StartClapTimeInput struct {
	BattleID uuid.UUID
}

type StartClapTimeUseCase struct {
}

func NewStartClapTimeUseCase() *StartClapTimeUseCase {
	return &StartClapTimeUseCase{}
}

func (uc *StartClapTimeUseCase) Execute(ctx context.Context, input StartClapTimeInput) error {
	fmt.Println("StartClapTimeUseCase Execute", input.BattleID)
	return nil
}
