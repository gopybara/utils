package redis

import "go.uber.org/fx"

func ProvideSentinel() fx.Option {
	return fx.Provide(
		NewRedisSentinelProvider,
	)
}

func ProvideRedis() fx.Option {
	return fx.Provide(
		NewRedisProvider,
	)
}

func AnnotateSentinel() fx.Option {
	return AnnotateRedis()
}

func AnnotateRedis() fx.Option {
	return fx.Provide(
		fx.Annotate(
			func(client Provider) Provider {
				return client
			},
			fx.ResultTags(`name:"redisClient"`),
		),
	)
}
