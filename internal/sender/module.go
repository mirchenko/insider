package sender

import "go.uber.org/fx"

var Module = fx.Provide(
	fx.Annotate(NewWorkerSender, fx.As(new(Sender))),
)
