package ratelimit

import (
	"context"
	"reflect"
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

type mockDs[T record] struct {
	record *T
}

func (m *mockDs[T]) Increment(ctx context.Context, key string) error {
	switch reflect.TypeOf(m.record).String() {
	case "*ratelimit.slidingWindowRecord":
		r := reflect.ValueOf(m.record).Interface().(*slidingWindowRecord)
		r.CurrentCount += 1
	}

	return nil
}

func (m *mockDs[T]) Get(ctx context.Context, key string) (record *T, err error) {
	return m.record, nil
}

func (m *mockDs[T]) Set(ctx context.Context, key string, record *T) error {
	m.record = record
	return nil
}
