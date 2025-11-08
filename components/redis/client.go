package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var (
	singleton redis.Cmdable
)

const (
	ModeCluster   = "cluster"
	ModeSingleton = "singleton"
	ModeSentinel  = "sentinel"
)

func InitRedis(conf *Config) error {
	var redisClient redis.Cmdable
	var err error
	password := conf.Password
	mode := conf.Mode
	if mode == "" {
		if conf.Cluster {
			mode = ModeCluster
		} else {
			mode = ModeSingleton
		}
	}
	switch mode {
	case ModeCluster:
		redisClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    conf.Addrs,
			Password: password,
		})
	case ModeSingleton:
		redisClient = redis.NewClient(&redis.Options{
			Addr:     conf.Addrs[0],
			Password: password,
		})
	case ModeSentinel:
		redisClient = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    conf.MasterName,
			SentinelAddrs: conf.Addrs,
			Password:      password,
		})
	}
	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	singleton = redisClient
	return nil
}

func GetRedisClient() redis.Cmdable {
	return singleton
}
