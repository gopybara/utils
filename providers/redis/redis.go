package redis

import (
	"github.com/go-redis/redis/v8"
)

type Provider interface {
	Client() *redis.Client
}

type redisProvider struct {
	client *redis.Client
}

func NewRedisProvider(config Config) (Provider, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.GetDSN(),
		Password: config.GetPassword(),
		DB:       config.GetDB(),
	})

	provider := &redisProvider{
		client: client,
	}

	return provider, nil
}

func NewRedisSentinelProvider(config SentinelConfig) (Provider, error) {
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:       config.GetMasterName(),
		SentinelAddrs:    config.GetNodes(),
		Password:         config.GetPassword(),
		SentinelPassword: config.GetSentinelPassword(),
		DB:               config.GetDB(),
	})

	provider := &redisProvider{
		client: client,
	}

	return provider, nil
}

func (p *redisProvider) Client() *redis.Client {
	return p.client
}
