package provider

import (
	"context"
	"insider/config"
	"insider/internal/model"
	"insider/pkg/logger"
	"time"

	"go.uber.org/fx"
	"resty.dev/v3"
)

type WebhookProviderResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

type WebhookProvider struct {
	logger *logger.Logger
	cfg    *config.Config
	client *resty.Client
}

func NewWebhookProvider(log *logger.Logger, cfg *config.Config, lc fx.Lifecycle) *WebhookProvider {
	client := resty.New().
		SetBaseURL(cfg.BaseURL).
		SetHeader("Content-Type", "application/json").
		SetRetryCount(cfg.RetriesCount).
		SetRetryWaitTime(time.Duration(cfg.RetryTimeoutSeconds)*time.Second).
		SetHeader("x-ins-auth-key", cfg.AuthKey).
		SetDebug(true)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Info().Msg("closing webhook provider http client")
			_ = client.Close()
			return nil
		},
	})

	return &WebhookProvider{
		logger: log,
		cfg:    cfg,
		client: client,
	}
}

func (p *WebhookProvider) Send(msg *model.Message) (*BaseProviderResponse, error) {
	res, err := p.client.R().
		SetBody(map[string]string{
			"to":      msg.PhoneNumber,
			"content": msg.Content,
		}).
		SetResult(&WebhookProviderResponse{}).
		Post(p.cfg.TenantID)

	if err != nil {
		p.logger.Error().Err(err).Msg("webhook provider error")
		return nil, err
	}

	response := res.Result().(*WebhookProviderResponse)

	return &BaseProviderResponse{
		MessageID: response.MessageID,
		Status:    response.Status,
	}, nil
}
