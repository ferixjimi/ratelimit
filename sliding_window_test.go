package ratelimit

import (
	"context"
	"testing"
	"time"
)

func TestSlidingWindowLimiter_Allow(t *testing.T) {
	limit := PerSecond(1)
	key := "key"
	ds := &mockDs[slidingWindowRecord]{record: &slidingWindowRecord{}}

	limiter := NewSlidingWindowLimiter(ds)

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

	if ds.record.CurrentCount != 1 || ds.record.PrevCount != 0 {
		t.Error("wrong current or prev Count")
	}

	time.Sleep(time.Second)

	result, err = limiter.Allow(context.Background(), key, limit)
	if err != nil {
		t.Error(err)
	}
	if !result.Allowed {
		t.Errorf("this attempt should be allowed")
	}

	if ds.record.CurrentCount != 1 || ds.record.PrevCount != 1 {
		t.Error("wrong current or prev Count")
	}
}
