package clapinfra

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/ABfry/album-battler/backend/internal/domain/service"
	"github.com/google/uuid"
)

var _ service.ClapWaitScheduler = (*ClapWaitScheduler)(nil)

type ClapWaitScheduler struct {
	mu           sync.Mutex
	triggered    map[uuid.UUID]map[uuid.UUID]bool
	defaultDelay time.Duration
}

func NewClapWaitScheduler(
	defaultDelay time.Duration,
) *ClapWaitScheduler {
	return &ClapWaitScheduler{
		triggered:    make(map[uuid.UUID]map[uuid.UUID]bool),
		defaultDelay: defaultDelay,
	}
}

func (s *ClapWaitScheduler) Schedule(battleID uuid.UUID, userID uuid.UUID, delay time.Duration, isEnd chan<- bool) {
	if delay <= 0 {
		delay = s.defaultDelay
	}

	s.mu.Lock()
	if s.triggered[battleID] == nil {
		s.triggered[battleID] = make(map[uuid.UUID]bool)
	} else if s.triggered[battleID][userID] {
		// すでにプレイヤーの拍手フェーズが終了している場合は何もしない
		s.mu.Unlock()
		return
	}

	ctx := context.Background()
	time.AfterFunc(delay, func() {
		if err := s.trigger(ctx, battleID, userID, isEnd); err != nil {
			log.Printf("failed to trigger clap wait timer for battle %s, user %s: %v", battleID, userID, err)
		}
	})
	s.mu.Unlock()
}

func (s *ClapWaitScheduler) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.triggered = map[uuid.UUID]map[uuid.UUID]bool{}
}

func (s *ClapWaitScheduler) trigger(ctx context.Context, battleID uuid.UUID, userID uuid.UUID, isEnd chan<- bool) error {
	s.mu.Lock()
	if s.triggered[battleID] == nil {
		s.triggered[battleID] = make(map[uuid.UUID]bool)
	} else if s.triggered[battleID][userID] {
		s.mu.Unlock()
		return nil
	}

	// 拍手フェーズが終了したことを記録
	s.triggered[battleID][userID] = true
	s.mu.Unlock()

	isEnd <- true

	// TODO: ResultStartのイベントに変更する

	return nil
}
