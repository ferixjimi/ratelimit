package ratelimit

import (
	"context"
	"time"
)

type FixedWindowLimiter struct {
	s Store[fixedWindowRecord]
}

type fixedWindowRecord struct {
	Start int64 `redis:"start"`
	Count int   `redis:"count"`
}

func NewFixedWindowLimiter(s Store[fixedWindowRecord]) Limiter {
	return &FixedWindowLimiter{s: s}
}

// todo: set error correctly
func (l *FixedWindowLimiter) Allow(ctx context.Context, key string, limit *Limit) (*Result, error) {
	if limit.Rate <= 0 {
		return &Result{
			Allowed: true,
		}, nil
	}

	bucket, err := l.s.Get(ctx, key)
	if err != nil || bucket.Start == 0 {
		l.s.Set(ctx, key, &fixedWindowRecord{
			Start: time.Now().UnixNano(),
			Count: 1,
		})

		return &Result{
			Allowed: true,
		}, nil
	}

	now := time.Now().UnixNano()
	windowLength := limit.Period.Nanoseconds()
	elapsedTime := now - bucket.Start
	if elapsedTime > windowLength {
		l.s.Set(ctx, key, &fixedWindowRecord{
			Start: now - (now-bucket.Start)%windowLength,
			Count: 1,
		})
		return &Result{
			Allowed: true,
		}, nil
	}

	if bucket.Count >= limit.Rate {
		return &Result{
			Allowed:    false,
			RetryAfter: time.Duration(l.retryAfter(elapsedTime, windowLength, now)),
		}, nil
	}

	l.s.Increment(ctx, key)
	return &Result{
		Allowed: true,
	}, nil
}

func (l *FixedWindowLimiter) retryAfter(elapsedTime, windowLength, now int64) int64 {
	return windowLength - (elapsedTime)%windowLength
}
