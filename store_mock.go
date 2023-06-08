package ratelimit

import (
	"context"
	"reflect"
)

type mockStore[T record] struct {
	record *T
}

func (m *mockStore[T]) Increment(ctx context.Context, key string) error {
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

func (m *mockStore[T]) Get(ctx context.Context, key string) (record *T, err error) {
	return m.record, nil
}

func (m *mockStore[T]) Set(ctx context.Context, key string, record *T) error {
	m.record = record
	return nil
}
