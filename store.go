package ratelimit

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type Datastore interface {
	Increment(ctx context.Context, key string) error

	Get(ctx context.Context, key string) (record *Record, err error)

	Set(ctx context.Context, key string, record *Record) error
}

const (
	currentCountFieldTag = "current"
)

type Record struct {
	Start        int64 `redis:"start"`
	PrevCount    int64 `redis:"prev"`
	CurrentCount int64 `redis:"current"`
}

func NewRedisStore(db *redis.Client) *RedisStore {
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

func (r *RedisStore) Get(ctx context.Context, key string) (*Record, error) {
	var record Record
	res := r.db.HGetAll(ctx, key)
	if err := res.Err(); err != nil {
		return nil, err
	}

	if err := res.Scan(&record); err != nil {
		return nil, res.Err()
	}

	return &record, nil
}

func (r *RedisStore) Set(ctx context.Context, key string, record *Record) error {
	return r.db.HSet(ctx, key, record).Err()
}
