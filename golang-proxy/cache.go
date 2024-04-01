package main

import (
	"context"

	"github.com/go-redis/redis"
)

type Cache interface {
	Save(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (string, error)
}

type RedisCache struct {
	Client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		Client: client,
	}
}

func (r *RedisCache) Save(ctx context.Context, key string, value string) error {
	return r.Client.Set(key, value, 0).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	result, err := r.Client.Get(key).Result()
	if err == redis.Nil {
		return "", nil // Cache miss
	}

	return result, err // Cache hit or error
}
