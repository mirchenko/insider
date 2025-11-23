package model

import (
	"strings"
	"time"
)

const (
	MessageStatusSent      = "sent"
	MessageStatusPending   = "pending"
	MessageStatusDelivered = "delivered"
	MessageStatusFailed    = "failed"
)

type Message struct {
	ID                int64      `json:"id" gorm:"primaryKey"`
	ExternalMessageID string     `json:"external_message_id"`
	Status            string     `json:"status"`
	PhoneNumber       string     `json:"phone_number"`
	Content           string     `json:"content"`
	Reason            *string    `json:"reason,omitempty"`
	CreatedAt         *time.Time `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

type MessageStatuses []string

func (s *MessageStatuses) UnmarshalParam(param string) error {
	if param == "" {
		return nil
	}

	*s = strings.Split(param, ",")
	return nil
}

type ListMessagesRequest struct {
	Limit  int             `form:"limit"`
	Offset int             `form:"offset"`
	Status MessageStatuses `form:"status" binding:"messageStatus"`
}

type ListMessagesResponse struct {
	Data                 []Message `json:"data"`
	ListResponseMetadata `json:"metadata"`
}
