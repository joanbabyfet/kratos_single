package data

import (
	"context"

	"kratos_single/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(c *conf.Data, logger log.Logger) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       int(c.Redis.Db),
	})

	// 测试连接
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	log.NewHelper(logger).Info("redis connected")

	return rdb
}