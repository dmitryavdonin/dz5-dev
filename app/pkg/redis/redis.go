package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
	ttl    time.Duration
}

func New(url, pass string, ttl time.Duration) (*Redis, error) {
	r := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: pass,
	})
	return &Redis{client: r, ttl: ttl}, nil
}

func (r *Redis) SetKey(ctx context.Context, key string, value interface{}) error {
	return r.client.Set(ctx, key, value, r.ttl).Err()
}

func (r *Redis) DeleteKey(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *Redis) GetKey(ctx context.Context, key string) (value interface{}, err error) {
	value, err = r.client.Get(ctx, key).Result()
	if err != nil {
		return
	}
	return
}

func (r *Redis) Close() error {
	return r.client.Close()
}
