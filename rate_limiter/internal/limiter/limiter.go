package limiter

import (
	"rate_limite/internal/storage"
	"time"
)

type RateLimiter struct {
	Storage storage.LimiterStorage
	DefaultLimit int
	DefaultWindow time.Duration
	DefaultBlock time.Duration
}

func NewRateLimiter(storage storage.LimiterStorage, limit int, window, block time.Duration) *RateLimiter {
	return &RateLimiter{
		Storage: storage,
		DefaultLimit: limit,
		DefaultWindow: window,
		DefaultBlock: block,
	}
}

func (rl *RateLimiter) Allow(key string, limit int, window, blockTime time.Duration) (bool, time.Time, error) {
	blocked, until, err := rl.Storage.IsBlocked(key)
	if err != nil {
		return false, time.Time{}, err
	}
	if blocked {
		return false, until, nil
	}

	count, err := rl.Storage.Increment(key, window)
	if err != nil {
		return false, time.Time{}, err
	}

	if count > limit {
		err := rl.Storage.Block(key, blockTime)
		if err != nil {
			return false, time.Time{}, err
		}
		return false, time.Now().Add(blockTime), nil
	}
	return true, time.Time{}, nil
}
