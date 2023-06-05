package ratelimit

import (
	"context"
	"github.com/go-redis/redismock/v9"
	"strconv"
	"testing"
	"time"
)

func TestRedisStore_Set(t *testing.T) {
	db, mock := redismock.NewClientMock()

	key := "key"
	record := Record{
		Start:        time.Now().UnixNano(),
		PrevCount:    1,
		CurrentCount: 2,
	}
	mock.ExpectHSet(key, record).SetVal(0)

	s := NewRedisStore(db)
	err := s.Set(context.Background(), key, &record)

	if err != nil {
		t.Error(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestRedisStore_Increment(t *testing.T) {
	db, mock := redismock.NewClientMock()

	key := "key"
	mock.ExpectHIncrBy(key, currentCountFieldTag, 1).SetVal(1)

	s := NewRedisStore(db)
	err := s.Increment(context.Background(), key)

	if err != nil {
		t.Error(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestRedisStore_Get(t *testing.T) {
	db, mock := redismock.NewClientMock()

	key := "key"
	record := Record{
		Start:        time.Now().UnixNano(),
		PrevCount:    1,
		CurrentCount: 2,
	}
	mock.ExpectHGetAll(key).SetVal(map[string]string{
		"start":   strconv.FormatInt(record.Start, 10),
		"prev":    "1",
		"current": "2",
	})

	s := NewRedisStore(db)
	result, err := s.Get(context.Background(), key)

	if err != nil {
		t.Error(err)
	}

	if *result != record {
		t.Errorf("expected: %+v, got: %+v", record, *result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
