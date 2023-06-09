package ratelimit

import (
	"context"
	"testing"
	"time"
)

func TestFixedWindowLimiter_Allow(t *testing.T) {
	limit := PerSecond(1)
	key := "key"
	s := &mockStore[fixedWindowRecord]{record: &fixedWindowRecord{}}

	limiter := NewFixedWindowLimiter(s)

	result, err := limiter.Allow(context.Background(), key, limit)
	if err != nil {
		t.Error(err)
	}

	if !result.Allowed {
		t.Error("this attempt should be allowed")
	}

	result, err = limiter.Allow(context.Background(), key, limit)
	if err != nil {
		t.Error(err)
	}

	if result.Allowed {
		t.Error("this attempt should'nt be allowed")
	}

	if s.record.Count != 1 {
		t.Error("wrong count")
	}

	if result.RetryAfter > time.Second {
		t.Error("retry after expected to be smaller than one second")
	}

}
