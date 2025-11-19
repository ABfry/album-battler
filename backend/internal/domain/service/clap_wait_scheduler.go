package service

import (
	"time"

	"github.com/google/uuid"
)

type ClapWaitScheduler interface {
	Schedule(battleID uuid.UUID, userID uuid.UUID, delay time.Duration, isEnd chan<- bool)
	Close()
}
