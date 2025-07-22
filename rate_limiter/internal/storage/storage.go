package storage

import "time"

type LimiterStorage interface {
	Increment(key string, window time.Duration) (int, error)
	GetCount(key string) (int, error)
	Block(key string, duration time.Duration) error
	IsBlocked(key string) (bool, time.Time, error)
	Reset(key string) error
}
