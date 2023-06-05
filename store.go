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

func NewRedisStore(host string, password string, db int) *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
	})

	return &RedisStore{
		client: rdb,
	}
}

type RedisStore struct {
	client *redis.Client
}

type Record struct {
	Start        int64 `redis:"start"`
	PrevCount    int64 `redis:"prev"`
	CurrentCount int64 `redis:"current"`
}

func (r *RedisStore) Increment(ctx context.Context, key string) error {
	return r.client.HIncrBy(ctx, key, "current", 1).Err()
}

func (r *RedisStore) Get(ctx context.Context, key string) (*Record, error) {
	var record Record
	res := r.client.HGetAll(ctx, key)
	if err := res.Err(); err != nil {
		return nil, err
	}

	if err := res.Scan(&record); err != nil {
		return nil, res.Err()
	}

	return &record, nil
}

func (r *RedisStore) Set(ctx context.Context, key string, record *Record) error {
	return r.client.HSet(ctx, key, record).Err()
}
