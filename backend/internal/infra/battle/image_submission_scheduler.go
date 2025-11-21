package battleinfra

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ABfry/album-battler/backend/internal/domain/service"
	"github.com/ABfry/album-battler/backend/internal/usecase/clap"
	"github.com/google/uuid"
)

var _ service.ImageSubmissionScheduler = (*ImageSubmissionScheduler)(nil)

type ImageSubmissionScheduler struct {
	mu           sync.Mutex
	timers       map[uuid.UUID]*time.Timer
	triggered    map[uuid.UUID]bool
	startClapUC  *clap.StartClapTimeUseCase
	defaultDelay time.Duration
}

func NewImageSubmissionScheduler(
	startClapUC *clap.StartClapTimeUseCase,
	defaultDelay time.Duration,
) *ImageSubmissionScheduler {
	return &ImageSubmissionScheduler{
		timers:       make(map[uuid.UUID]*time.Timer),
		triggered:    make(map[uuid.UUID]bool),
		startClapUC:  startClapUC,
		defaultDelay: defaultDelay,
	}
}

// 指定した遅延後に画像投稿を締め切り、拍手フェーズへ移行するタイマーを設定
func (s *ImageSubmissionScheduler) Schedule(battleID uuid.UUID, delay time.Duration) {
	if delay <= 0 {
		delay = s.defaultDelay
	}

	// すでに拍手フェーズに移行している場合は何もしない
	s.mu.Lock()
	if s.triggered[battleID] {
		s.mu.Unlock()
		return
	}

	// すでにタイマーが設定されている場合は停止
	if timer, ok := s.timers[battleID]; ok {
		timer.Stop()
	}

	// タイマーを設定
	ctx := context.Background()
	timer := time.AfterFunc(delay, func() {
		// 拍手フェーズへ移行
		fmt.Println("Move to clap phase by timer", battleID)
		if err := s.trigger(ctx, battleID); err != nil {
			log.Printf("failed to start clap time (timer) for battle %s: %v", battleID, err)
		}
	})
	s.timers[battleID] = timer
	s.mu.Unlock()
}

// タイマーを破棄し、即座に画像投稿を締め切り拍手フェーズへ移行
func (s *ImageSubmissionScheduler) TriggerNow(ctx context.Context, battleID uuid.UUID) error {
	fmt.Println("ClapScheduler TriggerNow", battleID)
	return s.trigger(ctx, battleID)
}

func (s *ImageSubmissionScheduler) Cancel(battleID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if timer, ok := s.timers[battleID]; ok {
		timer.Stop()
		delete(s.timers, battleID)
	}
}

func (s *ImageSubmissionScheduler) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for battleID, timer := range s.timers {
		timer.Stop()
		delete(s.timers, battleID)
	}
	s.triggered = map[uuid.UUID]bool{}
}

// 画像投稿を締め切り、拍手フェーズへ移行
func (s *ImageSubmissionScheduler) trigger(ctx context.Context, battleID uuid.UUID) error {
	s.mu.Lock()
	if s.triggered[battleID] {
		s.mu.Unlock()
		return nil
	}

	// タイマーを停止
	if timer, ok := s.timers[battleID]; ok {
		timer.Stop()
		delete(s.timers, battleID)
	}

	// 拍手フェーズに移行したことを記録
	s.triggered[battleID] = true
	s.mu.Unlock()

	// StartClapTimeUseCaseを実行
	if err := s.startClapUC.Execute(ctx, clap.StartClapTimeInput{BattleID: battleID}); err != nil {
		return err
	}
	return nil
}
