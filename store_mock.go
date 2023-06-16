package ratelimit

import (
	"context"
)

type mockStore struct {
	record interface{}
}

func (m *mockStore) Increment(ctx context.Context, key string) error {
	//switch any(m.record).(type) {
	//case *slidingWindowRecord:
	//	r := reflect.ValueOf(m.record).Interface().(*slidingWindowRecord)
	//	r.CurrentCount += 1
	//case *tokenBucketRecord:
	//	r := reflect.ValueOf(m.record).Interface().(*tokenBucketRecord)
	//	r.Count += 1
	//case *fixedWindowRecord:
	//	r := reflect.ValueOf(m.record).Interface().(*fixedWindowRecord)
	//	r.Count += 1
	//}

	return nil
}

func (m *mockStore) Get(ctx context.Context, key string) (record interface{}, err error) {
	return m.record, nil
}

func (m *mockStore) Set(ctx context.Context, key string, record interface{}) error {
	m.record = record
	return nil
}
