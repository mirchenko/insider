package services

import "go.uber.org/fx"

var Module = fx.Provide(
	fx.Annotate(NewMessagesService, fx.As(new(MessagesService))),
	fx.Annotate(NewSenderService, fx.As(new(SenderService))),
)
