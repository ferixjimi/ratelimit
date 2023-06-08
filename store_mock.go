package ratelimit

import (
	"context"
	"reflect"
)

type mockDs[T record] struct {
	record *T
}

func (m *mockDs[T]) Increment(ctx context.Context, key string) error {
	switch any(m.record).(type) {
	case *slidingWindowRecord:
		r := reflect.ValueOf(m.record).Interface().(*slidingWindowRecord)
		r.CurrentCount += 1
	case *tokenBucketRecord:
		r := reflect.ValueOf(m.record).Interface().(*tokenBucketRecord)
		r.Count += 1
	case *fixedWindowRecord:
		r := reflect.ValueOf(m.record).Interface().(*fixedWindowRecord)
		r.Count += 1
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
