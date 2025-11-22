package provider

import "insider/internal/model"

type BaseProviderResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

type Provider interface {
	Send(msg *model.Message) (*BaseProviderResponse, error)
}
