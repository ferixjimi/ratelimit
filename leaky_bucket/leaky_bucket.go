package leaky_bucket

import (
	"fmt"
	"github.com/ferixjimi/ratelimit"
	"time"
)

var _ ratelimit.Limiter = (*LeakyBucketLimiter)(nil)

// todo: add support for concurrent access
type LeakyBucketLimiter struct{}

func NewLeakyBucketLimiter() *LeakyBucketLimiter {
	return &LeakyBucketLimiter{}
}

type record struct {
	Last int64 `redis:"last"`
}

func (l *LeakyBucketLimiter) Allow(limit *ratelimit.Limit, inData interface{}) (allowed bool, outData interface{}, err error) {
	allowed = true
	if inData == nil {
		allowed = true
		outData = &record{
			Last: time.Now().UnixNano(),
		}
		return
	}

	bucket, ok := inData.(*record)
	if !ok {
		allowed = false
		err = fmt.Errorf("invalid in data type. want %T, got %T", &record{}, inData)
		return
	}

	elapsedTime := time.Now().UnixNano() - bucket.Last
	rate := limit.Period.Nanoseconds() / int64(limit.Rate)

	if elapsedTime < rate {
		allowed = false
	}

	bucket.Last = time.Now().UnixNano()
	outData = bucket
	return
}

func (l *LeakyBucketLimiter) retryAfter(elapsedTime, rate int64) int64 {
	return rate - elapsedTime
}
