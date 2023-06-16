package ratelimit

import (
	"time"
)

type Limiter interface { // todo add retry after
	Allow(limit *Limit, data interface{}) (bool, interface{}, error)
}

type Result struct {
	Allowed    bool
	RetryAfter time.Duration
}
