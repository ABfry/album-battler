package clapinfra

import (
	"sync"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
	"github.com/google/uuid"
)

var _ service.ClapCounter = (*MemoryClapCounter)(nil)

type MemoryClapCounter struct {
	mu     sync.RWMutex
	counts map[uuid.UUID]map[uuid.UUID]entity.Image
}

func NewMemoryClapCounter() *MemoryClapCounter {
	return &MemoryClapCounter{counts: make(map[uuid.UUID]map[uuid.UUID]entity.Image)}
}

func (c *MemoryClapCounter) Add(battleID, userID uuid.UUID, n int) entity.Image {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.counts[battleID] == nil {
		c.counts[battleID] = make(map[uuid.UUID]entity.Image)
	}

	img := c.counts[battleID][userID]
	img.AddUserScore(n)
	c.counts[battleID][userID] = img
	return img
}

func (c *MemoryClapCounter) Snapshot(battleID uuid.UUID) map[uuid.UUID]entity.Image {
	c.mu.RLock()
	defer c.mu.RUnlock()

	snapshot := make(map[uuid.UUID]entity.Image)
	for userID, count := range c.counts[battleID] {
		snapshot[userID] = count
	}
	return snapshot
}

func (c *MemoryClapCounter) Reset(battleID uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.counts, battleID)
}
