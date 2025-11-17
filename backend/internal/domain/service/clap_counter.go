package service

import (
	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/google/uuid"
)

type ClapCounter interface {
	Add(battleID, userID uuid.UUID, n int) entity.Image
	Snapshot(battleID uuid.UUID) map[uuid.UUID]entity.Image
	Reset(battleID uuid.UUID)
}
