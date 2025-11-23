package repository

import "go.uber.org/fx"

var Module = fx.Provide(
	fx.Annotate(NewMessagesRepository, fx.As(new(MessageRepository))),
)
