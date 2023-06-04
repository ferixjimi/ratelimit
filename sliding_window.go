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

	limited, err := l.ds.GetLimit(ctx, key)
	if err == nil && limited {
		return &Result{
			Allowed: false,
		}, nil
	}

	ttl := 2 * limit.Period.Milliseconds() / 1e3

	start, preCount, curCount, err := l.ds.Get(ctx, key, ttl)
	if err != nil {
		return nil, err
	}

	now := time.Now().UnixNano()
	unit := limit.Period.Nanoseconds()

	if (now - start) >= unit {
		start += unit
		preCount = curCount
		curCount = 0
		err = l.ds.Set(ctx, key, ttl, start, preCount, curCount)
		if err != nil {
			fmt.Println(err)
		}
	}

	d := float64(unit-(now-start)) / float64(unit)

	ec := float64(preCount)*d + float64(curCount)

	if ec >= float64(limit.Rate) {
		ttl := retryAfter(3, start, now, unit, preCount, curCount)
		l.ds.SetLimit(ctx, key, ttl/1e9)
		return &Result{
			Allowed:    false,
			RetryAfter: time.Duration(ttl),
		}, nil
	} else {
		err = l.ds.Add(ctx, key)
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
