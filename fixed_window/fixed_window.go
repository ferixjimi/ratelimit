package fixed_window

import (
	"fmt"
	"github.com/ferixjimi/ratelimit"
	"time"
)

var _ ratelimit.Limiter = (*FixedWindowLimiter)(nil)

type FixedWindowLimiter struct{}

func NewFixedWindowLimiter() *FixedWindowLimiter {
	return &FixedWindowLimiter{}
}

type record struct {
	Start int64 `redis:"start"`
	Count int   `redis:"count"`
}

func (l *FixedWindowLimiter) Allow(limit *ratelimit.Limit, inData interface{}) (allowed bool, outData interface{}, err error) {
	allowed = true
	if inData == nil {
		allowed = true
		outData = &record{
			Start: time.Now().UnixNano(),
			Count: 1,
		}
		return
	}

	window, ok := inData.(*record)
	if !ok {
		allowed = false
		err = fmt.Errorf("invalid in data type. want %T, got %T", &record{}, inData)
		return
	}

	now := time.Now().UnixNano()
	windowLength := limit.Period.Nanoseconds()
	elapsedTime := now - window.Start
	if elapsedTime > windowLength {
		allowed = true
		outData = &record{
			Start: now - (now-window.Start)%windowLength,
			Count: 1,
		}
		return
	}

	if window.Count >= limit.Rate {
		allowed = false
	}

	// increment
	window.Count += 1
	outData = window
	return
}

func (l *FixedWindowLimiter) retryAfter(elapsedTime, windowLength, now int64) int64 {
	return windowLength - (elapsedTime)%windowLength
}
