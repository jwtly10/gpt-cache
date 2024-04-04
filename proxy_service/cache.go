package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
)

const IDCounterKey = "global:nextItemId" // Key for the global ID counter

type Cache interface {
	Save(ctx context.Context, value []byte) (int64, error)
	Get(ctx context.Context, key int64) ([]byte, error)
}

type RedisCache struct {
	Client *redis.Client
}

func NewRedisClient(config RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: "",
		DB:       0,
	})
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		Client: client,
	}
}

// Save saves a new item and returns a unique ID for it.
func (r *RedisCache) Save(ctx context.Context, value []byte) (int64, error) {
	// Increment the global ID counter
	id, err := r.Client.Incr(IDCounterKey).Result()
	if err != nil {
		return 0, err
	}

	// Use the new ID to create a unique key for the item
	key := fmt.Sprintf("item:%d", id)

	// Save the item in Redis
	err = r.Client.Set(key, value, 0).Err()
	if err != nil {
		return 0, err
	}

	// Return the new ID
	return id, nil
}

func (r *RedisCache) Get(ctx context.Context, id int64) ([]byte, error) {
	// Convert the int64 ID to a string to form the key
	key := fmt.Sprintf("item:%d", id)

	// Use the key to retrieve the item from Redis
	// Note using .Bytes() as we are storing of the response raw bytes, not strings
	result, err := r.Client.Get(key).Bytes()
	if err == redis.Nil {
		return []byte{}, nil // Cache miss
	}
	return result, err // Cache hit or error
}
