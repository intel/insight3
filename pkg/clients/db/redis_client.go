package db

import (
	"context"

	redis "github.com/go-redis/redis/v8"
	"github.com/intel-sandbox/kube-score/pkg/common"
)

type RedisClient struct {
	Client *redis.Client
}

func (cli *RedisClient) NewClient(ctx context.Context, dbcondfig common.RunConfigDB) error {

	cli.Client = redis.NewClient(&redis.Options{
		Addr:     dbcondfig.Redis.Address,
		Password: dbcondfig.Redis.Password,
		DB:       dbcondfig.Redis.DB,
	})

	_, err := cli.Client.Ping(ctx).Result()
	return err
}
