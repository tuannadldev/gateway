package redis

import (
	"context"
	"errors"
	"gateway/config"
	"time"

	"github.com/redis/go-redis/v9"
)

// Init New Redis Client
func InitConnection(cfg *config.Config) (*redis.Client, error) {
	redisHost := cfg.Redis.RedisAddr

	if redisHost == "" {
		redisHost = ":6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:         redisHost,
		MinIdleConns: cfg.Redis.MinIdleConns,
		PoolSize:     cfg.Redis.PoolSize,
		PoolTimeout:  time.Duration(cfg.Redis.PoolTimeout) * time.Second,
		Password:     cfg.Redis.Password, // no password set
		DB:           cfg.Redis.DB,       // use default DB
	})
	if client.Ping(context.Background()).Err() != nil {
		return nil, errors.New("Can't connect redis server")
	}

	return client, nil
}
