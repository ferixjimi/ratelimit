package ratelimit

import (
	"context"
	"math"
	"time"
)

type TokenBucketLimiter struct {
	s Store[tokenBucketRecord]
}

type tokenBucketRecord struct {
	Start int64 `redis:"start"`
	Count int   `redis:"count"`
}

func NewTokenBucketLimiter(s Store[tokenBucketRecord]) Limiter {
	return &TokenBucketLimiter{s: s}
}

func (t *TokenBucketLimiter) Allow(ctx context.Context, key string, limit *Limit) (*Result, error) {
	if limit.Rate <= 0 {
		return &Result{
			Allowed: true,
		}, nil
	}

	bucket, err := t.s.Get(ctx, key)
	if err != nil || bucket.Start == 0 {
		t.s.Set(ctx, key, &tokenBucketRecord{
			Start: time.Now().UnixNano(),
			Count: limit.Rate - 1,
		})

		return &Result{
			Allowed: true,
		}, nil
	}

	newTokens := (time.Now().UnixNano() - bucket.Start) * int64(limit.Rate) / limit.Period.Nanoseconds()
	count := int(math.Min(float64(newTokens), float64(limit.Rate)))
	count = count + bucket.Count - 1

	if count < 0 {
		return &Result{
			Allowed:    false,
			RetryAfter: time.Duration(t.retryAfter(bucket, limit)),
		}, nil
	}

	bucket.Count = count
	bucket.Start += (newTokens * limit.Period.Nanoseconds()) / int64(limit.Rate)
	t.s.Set(ctx, key, bucket)

	return &Result{
		Allowed: true,
	}, nil
}

// todo: get now from args
func (t *TokenBucketLimiter) retryAfter(bucket *tokenBucketRecord, limit *Limit) int64 {
	speed := limit.Period.Nanoseconds() / int64(limit.Rate)
	return speed - (time.Now().UnixNano() - bucket.Start)
}
