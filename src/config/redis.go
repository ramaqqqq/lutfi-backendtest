package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"folkatech-customerIdentity/src/pkg/helpers"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type RedisCfg struct {
	Conn *redis.Client
}

type CacheConfig struct {
	RedisHost     string
	RedisPassword string
	RedisPort     string
	DSN           string
}

type Redis interface {
	WithCache(ctx context.Context, key string, dest interface{}, valFunc func() (interface{}, error)) error
	DelWithPattern(ctx context.Context, pattern string) error
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, duration time.Duration) error
	Del(ctx context.Context, key string) error
}

func LoadCacheConfig() CacheConfig {
	return CacheConfig{
		RedisHost:     viper.GetString("REDIS_HOST"),
		RedisPassword: viper.GetString("REDIS_PASSWORD"),
		RedisPort:     viper.GetString("REDIS_PORT"),
		DSN:           fmt.Sprintf("%s:%s", viper.GetString("REDIS_HOST"), viper.GetString("REDIS_PORT")),
	}
}

func CfgRedis(ctx context.Context, addr, password string) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	}

	redisClient := redis.NewClient(opts)
	if err := redisClient.Ping(ctx).Err(); err != nil {
		helpers.Logger("error", "init redis fail: "+err.Error())
		return nil, err
	}

	return redisClient, nil
}

func InitRedis(ctx context.Context, addr, password string) (Redis, error) {
	opts := &redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	}

	redisClient := redis.NewClient(opts)
	if err := redisClient.Ping(ctx).Err(); err != nil {
		helpers.Logger("error", "init redis fail: "+err.Error())
		return nil, err
	}

	rediss := RedisCfg{
		Conn: redisClient,
	}

	return &rediss, nil
}

func (rds *RedisCfg) WithCache(ctx context.Context, key string, dest interface{}, valFunc func() (interface{}, error)) error {
	val, err := rds.Conn.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		helpers.Logger("error", "error when get data redis: "+err.Error())
	}

	if val != "" {
		err := json.Unmarshal([]byte(val), dest)
		if err == nil {
			return nil
		}

		helpers.Logger("error", "error when unmarshal data redis: "+err.Error())
	}

	data, err := valFunc()
	if err != nil {
		helpers.Logger("error", "error function params:  "+err.Error())
		return err
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		helpers.Logger("error", "error when marshal dataJSON for redis:  "+err.Error())
		return err
	}

	err = rds.Conn.Set(ctx, key, dataJSON, 0).Err()
	if err != nil {
		helpers.Logger("error", "error when set data redis: "+err.Error())

	}

	err = json.Unmarshal(dataJSON, dest)
	if err != nil {
		helpers.Logger("error", "error when unmarshal data for return: "+err.Error())
		return err
	}

	return nil
}

func (rds *RedisCfg) Get(ctx context.Context, key string) (string, error) {
	val, err := rds.Conn.Get(ctx, key).Result()
	if err != nil {
		helpers.Logger("error", "error when get data redis: "+err.Error())

		return "", err
	}

	return val, nil
}

func (rds *RedisCfg) Set(ctx context.Context, key string, value string, duration time.Duration) error {

	err := rds.Conn.Set(ctx, key, value, duration).Err()
	if err != nil {
		helpers.Logger("error", "error when set data redis: "+err.Error())
		return err
	}

	return nil
}

func (rds *RedisCfg) Del(ctx context.Context, key string) error {
	err := rds.Conn.Del(ctx, key).Err()
	if err != nil {
		helpers.Logger("error", "error when delete data redis:  "+err.Error())
		return err
	}

	return nil
}

func (rds *RedisCfg) DelWithPattern(ctx context.Context, pattern string) error {
	var cursor uint64
	var keys []string

	for {
		var err error
		keys, cursor, err = rds.Conn.Scan(ctx, cursor, pattern, 1000).Result()
		if err != nil {
			helpers.Logger("error", "something wrong when scan keys: %v"+err.Error())
			return err
		}

		if len(keys) == 0 {
			break
		}

		err = rds.Conn.Del(ctx, keys...).Err()
		if err != nil {
			helpers.Logger("error", "something wrong when deleted data: %v"+err.Error())
			return err
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}
