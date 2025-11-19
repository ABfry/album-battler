package service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// 画像投稿期限を管理し、期限後に拍手フェーズへの遷移を調整
type ImageSubmissionScheduler interface {
	// 指定した遅延後に画像投稿期限を迎え、拍手フェーズへ移行するタイマーを設定する。
	Schedule(battleID uuid.UUID, delay time.Duration)
	// タイマーを破棄し、即座に画像投稿を締め切って拍手フェーズへ移行させる。
	TriggerNow(ctx context.Context, battleID uuid.UUID) error
	// 拍手フェーズに移行せずにスケジュール済みのタイマーをキャンセルする。
	Cancel(battleID uuid.UUID)
	// すべてのタイマーを停止し、内部リソースを解放する。
	Close()
}
