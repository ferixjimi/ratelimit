package ratelimit

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type Datastore[T record] interface {
	Increment(ctx context.Context, key string) error

	Get(ctx context.Context, key string) (record *T, err error)

	Set(ctx context.Context, key string, record *T) error
}

type record interface {
	slidingWindowRecord | tokenBucketRecord | leakyBucketRecord | fixedWindowRecord
}

func NewRedisStore[T record](db *redis.Client) Datastore[T] {
	return &RedisStore[T]{
		db: db,
	}
}

type RedisStore[T record] struct {
	db *redis.Client
}

func (r *RedisStore[T]) Increment(ctx context.Context, key string) error {
	return r.db.HIncrBy(ctx, key, currentCountFieldTag, 1).Err()
}

func (r *RedisStore[T]) Get(ctx context.Context, key string) (*T, error) {
	res := r.db.HGetAll(ctx, key)
	if err := res.Err(); err != nil {
		return nil, err
	}

	var record T
	if err := res.Scan(&record); err != nil {
		return nil, err
	}

	return &record, nil
}

func (r *RedisStore[T]) Set(ctx context.Context, key string, record *T) error {
	return r.db.HSet(ctx, key, record).Err()
}
