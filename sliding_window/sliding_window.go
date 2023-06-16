package sliding_window

import (
	"fmt"
	"github.com/ferixjimi/ratelimit"
	"time"
)

var _ ratelimit.ILimiter = (*SlidingWindowLimiter)(nil)

type SlidingWindowLimiter struct{}

func NewSlidingWindowLimiter() *SlidingWindowLimiter {
	return &SlidingWindowLimiter{}
}

type record struct {
	Start        int64 `redis:"start"`
	PrevCount    int64 `redis:"prev"`
	CurrentCount int64 `redis:"current"`
}

func (l *SlidingWindowLimiter) Allow(limit *ratelimit.Limit, inData interface{}) (allowed bool, outData interface{}, err error) {
	allowed = true
	if inData == nil {
		allowed = true
		outData = &record{
			Start:        time.Now().UnixNano(),
			CurrentCount: 1,
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

	if (now - window.Start) >= windowLength {
		window.Start += windowLength
		window.PrevCount = window.CurrentCount
		window.CurrentCount = 0
	}

	d := float64(windowLength-(now-window.Start)) / float64(windowLength)

	currentCount := float64(window.PrevCount)*d + float64(window.CurrentCount)

	if currentCount >= float64(limit.Rate) {
		allowed = false
	}

	// increment
	window.CurrentCount += 1
	return
}

func (l *SlidingWindowLimiter) retryAfter(size, start, now, unit int64, preCount int64, curCount int64) int64 {
	d := 1.
	if preCount != 0 {
		d -= float64(size-curCount) / float64(preCount)
	}
	x := d*float64(unit) + float64(start)
	return int64(x) - now
}
