package handler

import "go.uber.org/fx"

var Module = fx.Provide(NewMessagesHandler, NewSenderHandler)
