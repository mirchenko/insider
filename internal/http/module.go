package http

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewRouters),
	fx.Provide(NewServer),
	fx.Invoke(StartServer),
)
