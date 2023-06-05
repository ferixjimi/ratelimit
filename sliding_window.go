package ratelimit

import (
	"context"
	"fmt"
	"time"
)

type SlidingWindowLimiter struct {
	ds Datastore
}

func NewSlidingWindowLimiter(ds Datastore) Limiter {
	return &SlidingWindowLimiter{
		ds: ds,
	}
}

func (l *SlidingWindowLimiter) Allow(ctx context.Context, key string, limit Limit) (*Result, error) {
	// ignore negative rates
	if limit.Rate < 0 {
		return &Result{
			Allowed: true,
		}, nil
	}

	window, err := l.ds.Get(ctx, key)
	if err != nil || window.Start == 0 {
		l.ds.Set(ctx, key, &Record{
			Start:        time.Now().UnixNano(),
			CurrentCount: 1,
		})

		return &Result{
			Allowed: true,
		}, nil
	}

	now := time.Now().UnixNano()
	windowLength := limit.Period.Nanoseconds()

	if (now - window.Start) >= windowLength {
		window.Start += windowLength
		window.PrevCount = window.CurrentCount
		window.CurrentCount = 0
		err = l.ds.Set(ctx, key, window)
		if err != nil {
			fmt.Println(err)
		}
	}

	d := float64(windowLength-(now-window.Start)) / float64(windowLength)

	currentCount := float64(window.PrevCount)*d + float64(window.CurrentCount)

	if currentCount >= float64(limit.Rate) {
		ttl := retryAfter(int64(limit.Rate), window.Start, now, windowLength, window.PrevCount, window.CurrentCount)
		return &Result{
			Allowed:    false,
			RetryAfter: time.Duration(ttl),
		}, nil
	} else {
		err = l.ds.Increment(ctx, key)
		if err != nil {
			fmt.Println(err)
		}
		return &Result{
			Allowed: true,
		}, nil
	}
}

func retryAfter(size, start, now, unit int64, preCount int64, curCount int64) int64 {
	d := 1.
	if preCount != 0 {
		d -= float64(size-curCount) / float64(preCount)
	}
	x := d*float64(unit) + float64(start)
	return int64(x) - now
}

type Result struct {
	Limit      Limit
	Allowed    bool
	RetryAfter time.Duration
}
