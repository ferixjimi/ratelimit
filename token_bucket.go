package ratelimit

import (
	"context"
	"math"
	"time"
)

type TokenBucketLimiter struct {
	ds Datastore[tokenBucketRecord]
}

type tokenBucketRecord struct {
	Start int64 `redis:"start"`
	Count int   `redis:"count"`
}

func NewTokenBucketLimiter(ds Datastore[tokenBucketRecord]) Limiter {
	return &TokenBucketLimiter{ds: ds}
}

func (t *TokenBucketLimiter) Allow(ctx context.Context, key string, limit *Limit) (*Result, error) {
	if limit.Rate < 0 {
		return &Result{
			Allowed: true,
		}, nil
	}

	bucket, err := t.ds.Get(ctx, key)
	if err != nil || bucket.Start == 0 {
		t.ds.Set(ctx, key, &tokenBucketRecord{
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
	t.ds.Set(ctx, key, bucket)

	return &Result{
		Allowed: true,
	}, nil
}

func (t *TokenBucketLimiter) retryAfter(bucket *tokenBucketRecord, limit *Limit) int64 {
	speed := limit.Period.Nanoseconds() / int64(limit.Rate)
	return speed - (time.Now().UnixNano() - bucket.Start)
}
