package ratelimit

import (
	"context"
	"testing"
	"time"
)

func TestLeakyBucketLimiter_Allow(t *testing.T) {
	limit := PerSecond(1)
	key := "key"
	ds := &mockDs[leakyBucketRecord]{record: &leakyBucketRecord{}}

	limiter := NewLeakyBucketLimiter(ds)

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

	if result.RetryAfter > time.Second {
		t.Error("retry after expected to smaller than one second")
	}

}
