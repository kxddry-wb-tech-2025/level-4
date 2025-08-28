package redis

import (
	"calendar/internal/config"
	"calendar/internal/storage"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Storage is the Redis storage client
type Storage struct {
	client *redis.Client
}

// NewStorage creates a new Redis storage client
func NewStorage(config *config.RedisConfig) *Storage {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Username: config.Username,
		Password: config.Password,
		DB:       config.Database,
	})

	return &Storage{client: client}
}

const (
	keyDueZset = "notify:due"
)

// PopDue returns the due elements from a sorted set with the given key
func (s *Storage) PopDue(ctx context.Context, key string, limit int64) ([]string, error) {
	vals, err := s.client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    fmt.Sprintf("%d", time.Now().Unix()),
		Offset: 0,
		Count:  limit,
	}).Result()
	if err != nil {
		return nil, err
	}

	if len(vals) == 0 {
		return nil, nil
	}

	if err := s.client.ZRem(ctx, key, anySlice(vals)...).Err(); err != nil {
		return nil, err
	}

	return vals, nil
}

func anySlice[T any](in []T) []any {
	out := make([]any, 0, len(in))
	for _, v := range in {
		out = append(out, v)
	}
	return out
}

// Enqueue adds a notification id to the delayed queue
func (s *Storage) Enqueue(ctx context.Context, id string, at time.Time) error {
	return s.client.ZAdd(ctx, keyDueZset, redis.Z{
		Score:  float64(at.Unix()),
		Member: id,
	}).Err()
}

// Remove deletes a scheduled notification id from the delayed queue
func (s *Storage) Remove(ctx context.Context, id string) error {
	err := s.client.ZRem(ctx, keyDueZset, id).Err()
	if errors.Is(err, redis.Nil) {
		return storage.ErrNotFound
	}
	return err
}
