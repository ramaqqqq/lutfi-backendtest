package db

import (
	"context"
	"fmt"
	"folkatech-customerIdentity/src/config"

	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg *config.Config) *redis.Client {
	redisAddr := fmt.Sprintf("%s:%s", cfg.Cache.RedisHost, cfg.Cache.RedisPort)
	client, err := config.CfgRedis(context.Background(), redisAddr, cfg.Cache.RedisPassword)
	if err != nil {
		panic(err)
	}
	return client
}
