package event

import "time"

type DomainEvent interface {
	// イベントのタイプを返す
	EventType() string

	// イベントの発生時刻を返す
	OccurredAt() time.Time
}

// ドメインイベントを記録するための構造体
type AggregateRoot struct {
	events []DomainEvent
}

// ドメインイベントを記録
func (a *AggregateRoot) RecordEvent(event DomainEvent) {
	a.events = append(a.events, event)
}

// 記録されたイベントを取得し、クリア
func (a *AggregateRoot) PopEvents() []DomainEvent {
	events := a.events
	a.events = []DomainEvent{}
	return events
}

// 未処理のイベントがあるか
func (a *AggregateRoot) HasEvents() bool {
	return len(a.events) > 0
}
