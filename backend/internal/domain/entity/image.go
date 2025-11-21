package entity

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Image struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	BattleID      uuid.UUID
	ImageURL      string
	UploadedAt    time.Time
	AIScore       float64 // 0.0 ~ 100.0
	UserScore     int     // 0 ~ N
	AIExplanation string  // AIによる評価の説明文
}

func NewImage(userID, battleID uuid.UUID, imageURL string) (*Image, error) {
	if userID == uuid.Nil {
		return nil, errors.New("userID is required")
	}
	if battleID == uuid.Nil {
		return nil, errors.New("battleID is required")
	}
	if imageURL == "" {
		return nil, errors.New("imageURL is required")
	}

	return &Image{
		ID:         uuid.New(),
		UserID:     userID,
		BattleID:   battleID,
		ImageURL:   imageURL,
		UploadedAt: time.Now(),
	}, nil
}

func (i *Image) SetScore(aiScore float64, userScore int) error {
	if aiScore < 0 || aiScore > 100 {
		return errors.New("ai score must be between 0 and 100")
	}
	if userScore < 0 {
		return errors.New("user score must be non-negative")
	}

	i.AIScore = aiScore
	i.UserScore = userScore
	return nil
}

func (i *Image) AddUserScore(score int) error {
	if score < 0 || score > 10 {
		fmt.Printf("invalid score to add: %d", score)
		return errors.New("score to add must be between 0 and 10")
	}

	if i.UserScore+score > 10 {
		fmt.Printf("user score exceeds maximum after adding: current=%d, adding=%d", i.UserScore, score)
		return nil
	}
	i.UserScore += score
	return nil
}

func (i *Image) SetAIExplanation(explanation string) {
	i.AIExplanation = explanation
}
