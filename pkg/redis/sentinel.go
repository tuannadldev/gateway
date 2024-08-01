package redis

import (
	"context"
	"errors"
	"gateway/config"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// Init NewFailoverClusterClient
func InitRedis(cfg *config.Config) (*redis.ClusterClient, error) {
	addrs := strings.Split(cfg.RedisSentinel.Addr, ",")

	rdb := redis.NewFailoverClusterClient(&redis.FailoverOptions{
		MasterName:       cfg.RedisSentinel.MasterName,
		SentinelAddrs:    addrs,
		DB:               0,
		Password:         cfg.RedisSentinel.Password,
		SentinelPassword: cfg.RedisSentinel.SentinelPassword,
		PoolSize:         cfg.RedisSentinel.PoolSize,
		MinIdleConns:     cfg.RedisSentinel.MinIdleConns,
		MaxRetries:       cfg.RedisSentinel.MaxRetries,
		MinRetryBackoff:  cfg.RedisSentinel.MinRetryBackoff * time.Millisecond,
		MaxRetryBackoff:  cfg.RedisSentinel.MaxRetryBackoff * time.Millisecond,
		DialTimeout:      cfg.RedisSentinel.DialTimeout * time.Millisecond,
		ReadTimeout:      cfg.RedisSentinel.ReadTimeout * time.Millisecond,
		WriteTimeout:     cfg.RedisSentinel.WriteTimeout * time.Millisecond,
		PoolFIFO:         cfg.RedisSentinel.PoolFIFO,
		PoolTimeout:      cfg.RedisSentinel.PoolTimeout * time.Millisecond,
		// RouteByLatency:   true, //
		// RouteRandomly:    true, //
	})

	if rdb.Ping(context.Background()).Err() != nil {
		return nil, errors.New("Can't connect redis sentinel")
	}
	return rdb, nil
}
