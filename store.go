package ratelimit

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Datastore interface {
	Add(ctx context.Context, key string) error

	Get(ctx context.Context, key string, ttl int64) (start int64, preCount int64, curCount int64, err error)

	Set(ctx context.Context, key string, ttl int64, start int64, preCount int64, curCount int64) error

	GetLimit(ctx context.Context, key string) (bool, error)
	SetLimit(ctx context.Context, key string, ttl int64) error
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
	Ttl          int64 `redis:"ttl"` // fixme
}

func (r *RedisStore) Add(ctx context.Context, key string) error {
	record := Record{
		Start: time.Now().Unix(),
	}
	return r.client.HSet(ctx, key, record).Err()
}

func (r *RedisStore) Get(ctx context.Context, key string, ttl int64) (start int64, preCount int64, curCount int64, err error) {
	var record Record
	err = r.client.HGetAll(ctx, key).Scan(&record)
	if err != nil {
		return 0, 0, 0, err
	}
	return record.Start, record.PrevCount, record.CurrentCount, nil
}

func (r *RedisStore) Set(ctx context.Context, key string, ttl int64, start int64, preCount int64, curCount int64) error {
	record := Record{
		Start:        start,
		PrevCount:    preCount,
		CurrentCount: curCount,
		Ttl:          ttl,
	}
	return r.client.HSet(ctx, key, record).Err()
}

func (r *RedisStore) GetLimit(ctx context.Context, key string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisStore) SetLimit(ctx context.Context, key string, ttl int64) error {
	//TODO implement me
	panic("implement me")
}
