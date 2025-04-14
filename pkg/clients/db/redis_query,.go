package db

import (
	"context"
	"fmt"

	redis "github.com/go-redis/redis/v8"
)

func (rdb *RedisClient) AddKey(ctx context.Context, key, value string) error {

	if err := rdb.Client.Set(ctx, key, value, 0).Err(); err != nil {
		return fmt.Errorf("error insering key: %s", key)
	}
	return nil
}

func (rdb *RedisClient) GetVaule(ctx context.Context, key string) (string, error) {

	val, err := rdb.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key not found: %s", key)
	} else if err != nil {
		return "", fmt.Errorf("error quering database")
	}
	return val, nil
}
