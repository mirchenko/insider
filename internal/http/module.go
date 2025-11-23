package http

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewHealthHandler),
	fx.Provide(NewMessagesHandler),
	fx.Provide(NewSenderHandler),
	fx.Provide(NewRouters),
	fx.Provide(NewServer),
	fx.Invoke(StartServer),
)
