package token_bucket

import (
	"fmt"
	"github.com/ferixjimi/ratelimit"
	"math"
	"time"
)

type TokenBucketLimiter struct{}

func NewTokenBucketLimiter() *TokenBucketLimiter {
	return &TokenBucketLimiter{}
}

type record struct {
	Start int64 `redis:"start"`
	Count int   `redis:"count"`
}

func (t *TokenBucketLimiter) Allow(limit *ratelimit.Limit, inData interface{}) (allowed bool, outData interface{}, err error) {
	allowed = true
	if inData == nil {
		allowed = true
		outData = &record{
			Start: time.Now().UnixNano(),
			Count: limit.Rate - 1,
		}
		return
	}

	bucket, ok := inData.(*record)
	if !ok {
		allowed = false
		err = fmt.Errorf("invalid in data type. want %T, got %T", &record{}, inData)
		return
	}

	newTokens := (time.Now().UnixNano() - bucket.Start) * int64(limit.Rate) / limit.Period.Nanoseconds()
	count := int(math.Min(float64(newTokens), float64(limit.Rate)))
	count = count + bucket.Count - 1

	if count < 0 {
		allowed = false
	}

	bucket.Count = count
	bucket.Start += (newTokens * limit.Period.Nanoseconds()) / int64(limit.Rate)

	outData = bucket
	return
}

// todo: get now from args
func (t *TokenBucketLimiter) retryAfter(bucket *record, limit *ratelimit.Limit) int64 {
	speed := limit.Period.Nanoseconds() / int64(limit.Rate)
	return speed - (time.Now().UnixNano() - bucket.Start)
}
