package ratelimit

import (
	"context"
	"reflect"
	"sync"
)

type LocalStore[T record] struct {
	mu    sync.RWMutex
	store map[string]*T
}

func NewLocalStore[T record]() *LocalStore[T] {
	return &LocalStore[T]{
		store: map[string]*T{},
	}
}

func (ls *LocalStore[T]) Increment(ctx context.Context, key string) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	val := ls.store[key]
	switch any(val).(type) {
	case *slidingWindowRecord:
		r := reflect.ValueOf(val).Interface().(*slidingWindowRecord)
		r.CurrentCount += 1
	case *tokenBucketRecord:
		r := reflect.ValueOf(val).Interface().(*tokenBucketRecord)
		r.Count += 1
	case *fixedWindowRecord:
		r := reflect.ValueOf(val).Interface().(*fixedWindowRecord)
		r.Count += 1
	}

	return nil
}

func (ls *LocalStore[T]) Get(ctx context.Context, key string) (record *T, err error) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()
	return ls.store[key], nil
}

func (ls *LocalStore[T]) Set(ctx context.Context, key string, record *T) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	ls.store[key] = record
	return nil
}
