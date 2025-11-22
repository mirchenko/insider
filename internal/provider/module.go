package provider

import "go.uber.org/fx"

var Module = fx.Provide(NewWebhookProvider)
