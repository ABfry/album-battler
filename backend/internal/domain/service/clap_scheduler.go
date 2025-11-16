package service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// バトルの投稿締切を管理し、拍手フェーズへの遷移を調整
type ClapScheduler interface {
	// 指定した遅延後に拍手フェーズへ移行するタイマーを設定する。
	Schedule(battleID uuid.UUID, delay time.Duration)
	// タイマーを破棄し、即座に拍手フェーズへ移行させる。
	TriggerNow(ctx context.Context, battleID uuid.UUID) error
	// 拍手フェーズに移行せずにスケジュール済みのタイマーをキャンセルする。
	Cancel(battleID uuid.UUID)
	// すべてのタイマーを停止し、内部リソースを解放する。
	Close()
}
