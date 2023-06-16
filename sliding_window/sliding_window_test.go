package sliding_window

import (
	"context"
	"testing"
	"time"
)

func TestSlidingWindowLimiter_Allow(t *testing.T) {
	limit := PerSecond(1)
	key := "key"
	s := &mockStore[slidingWindowRecord]{record: &slidingWindowRecord{}}

	limiter := NewSlidingWindowLimiter(s)

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

	if s.record.CurrentCount != 1 || s.record.PrevCount != 0 {
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

	if s.record.CurrentCount != 1 || s.record.PrevCount != 1 {
		t.Error("wrong current or prev Count")
	}
}
