package ratelimit

import (
	"context"
	"time"
)

// todo: add support for concurrent access
type LeakyBucketLimiter struct {
	s Store[leakyBucketRecord]
}

type leakyBucketRecord struct {
	Last int64 `redis:"last"`
}

func NewLeakyBucketLimiter(s Store[leakyBucketRecord]) Limiter {
	return &LeakyBucketLimiter{s: s}
}

func (l *LeakyBucketLimiter) Allow(ctx context.Context, key string, limit *Limit) (*Result, error) {
	if limit.Rate <= 0 {
		return &Result{
			Allowed: true,
		}, nil
	}

	bucket, err := l.s.Get(ctx, key)
	if err != nil || bucket.Last == 0 {
		l.s.Set(ctx, key, &leakyBucketRecord{
			Last: time.Now().UnixNano(),
		})

		return &Result{
			Allowed: true,
		}, nil
	}

	elapsedTime := time.Now().UnixNano() - bucket.Last
	rate := limit.Period.Nanoseconds() / int64(limit.Rate)

	if elapsedTime < rate {
		return &Result{
			Allowed:    false,
			RetryAfter: time.Duration(l.retryAfter(elapsedTime, rate)),
		}, nil
	}

	bucket.Last = time.Now().UnixNano()
	l.s.Set(ctx, key, bucket)

	return &Result{
		Allowed: true,
	}, nil
}

func (l *LeakyBucketLimiter) retryAfter(elapsedTime, rate int64) int64 {
	return rate - elapsedTime
}
