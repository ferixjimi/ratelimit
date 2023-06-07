package ratelimit

import (
	"time"
)

type Limit struct {
	Rate   int
	Period time.Duration
}

func PerSecond(rate int) *Limit {
	return &Limit{
		Rate:   rate,
		Period: time.Second,
	}
}

func PerMinute(rate int) *Limit {
	return &Limit{
		Rate:   rate,
		Period: time.Minute,
	}
}

func PerHour(rate int) *Limit {
	return &Limit{
		Rate:   rate,
		Period: time.Hour,
	}
}
