package ratelimit

import (
	"context"
	"sync"
)

type LocalStore struct {
	mu    sync.RWMutex
	store map[string]interface{}
}

func NewLocalStore() *LocalStore {
	return &LocalStore{
		store: map[string]interface{}{},
	}
}

func (ls *LocalStore) Increment(ctx context.Context, key string) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	//val := ls.store[key]
	//switch val.(type) {
	//case *slidingWindowRecord:
	//	r := reflect.ValueOf(val).Interface().(*slidingWindowRecord)
	//	r.CurrentCount += 1
	//case *tokenBucketRecord:
	//	r := reflect.ValueOf(val).Interface().(*tokenBucketRecord)
	//	r.Count += 1
	//case *fixedWindowRecord:
	//	r := reflect.ValueOf(val).Interface().(*fixedWindowRecord)
	//	r.Count += 1
	//}

	return nil
}

func (ls *LocalStore) Get(ctx context.Context, key string) (record interface{}, err error) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()
	return ls.store[key], nil
}

func (ls *LocalStore) Set(ctx context.Context, key string, record interface{}) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	ls.store[key] = record
	return nil
}
