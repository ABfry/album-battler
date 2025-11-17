package clapinfra

import (
	"sync"

	"github.com/google/uuid"
)

type MemoryClapCounter struct {
	mu     sync.RWMutex
	counts map[uuid.UUID]map[uuid.UUID]int
}

func NewMemoryClapCounter() *MemoryClapCounter {
	return &MemoryClapCounter{counts: make(map[uuid.UUID]map[uuid.UUID]int)}
}

func (c *MemoryClapCounter) Add(battleID, userID uuid.UUID, n int) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.counts[battleID] == nil {
		c.counts[battleID] = make(map[uuid.UUID]int)
	}
	c.counts[battleID][userID] += n
	return c.counts[battleID][userID]
}

func (c *MemoryClapCounter) Snapshot(battleID uuid.UUID) map[uuid.UUID]int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	snapshot := make(map[uuid.UUID]int)
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
