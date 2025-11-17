package service

import "github.com/google/uuid"

type ClapCounter interface {
	Add(battleID, userID uuid.UUID, n int) int
	Snapshot(battleID uuid.UUID) map[uuid.UUID]int
	Reset(battleID uuid.UUID)
}
