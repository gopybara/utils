package supervisor

import "go.uber.org/fx"

func ProvideGoroutineSupervisor() fx.Option {
	return fx.Provide(
		NewGoroutineSupervisor,
	)
}
