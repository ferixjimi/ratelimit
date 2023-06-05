package ratelimit

import (
	"context"
	"testing"
	"time"
)

func TestSlidingWindowLimiter_Allow(t *testing.T) {
	limit := PerSecond(1)
	key := "key"
	ds := &mockDs{record: &Record{}}

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
		t.Error("wrong current or prev count")
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
		t.Error("wrong current or prev count")
	}
}

type mockDs struct {
	record *Record
}

func (m *mockDs) Increment(ctx context.Context, key string) error {
	m.record.CurrentCount += 1
	return nil
}

func (m *mockDs) Get(ctx context.Context, key string) (record *Record, err error) {
	return m.record, nil
}

func (m *mockDs) Set(ctx context.Context, key string, record *Record) error {
	m.record = record
	return nil
}
