package app

import (
	"context"
	"fmt"
	"userservice/internal/config"

	"github.com/redis/go-redis/v9"
)

func mustLoadRedis(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisConf.Host, cfg.RedisConf.Port),
		Password: cfg.RedisConf.Password,
		DB:       cfg.RedisConf.DB,
	})

	if err := client.Ping(context.Background()); err.Err() != nil {
		panic("failed to connect to the redis: " + err.Err().Error())
	}

	return client
}
