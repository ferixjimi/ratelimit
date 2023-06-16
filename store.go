package ratelimit

import (
	"context"
	"github.com/redis/go-redis/v9"
)

const (
	currentCountFieldTag = "current"
)

type Store interface {
	Increment(ctx context.Context, key string) error

	Get(ctx context.Context, key string) (record interface{}, err error)

	Set(ctx context.Context, key string, record interface{}) error
}

func NewRedisStore(db *redis.Client) Store {
	return &RedisStore{
		db: db,
	}
}

type RedisStore struct {
	db *redis.Client
}

func (r *RedisStore) Increment(ctx context.Context, key string) error {
	return r.db.HIncrBy(ctx, key, currentCountFieldTag, 1).Err()
}

func (r *RedisStore) Get(ctx context.Context, key string) (interface{}, error) {
	res := r.db.HGetAll(ctx, key)
	if err := res.Err(); err != nil {
		return nil, err
	}

	var record interface{}
	if err := res.Scan(&record); err != nil {
		return nil, err
	}

	return &record, nil
}

func (r *RedisStore) Set(ctx context.Context, key string, record interface{}) error {
	return r.db.HSet(ctx, key, record).Err()
}
