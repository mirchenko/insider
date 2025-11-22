package model

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

type ListMessagesRequest struct {
	Limit  int      `form:"limit"`
	Offset int      `form:"offset"`
	Status []string `form:"status" binding:"messageStatus"`
}

type ListMessagesResponse struct {
	Data                 []Message `json:"data"`
	ListResponseMetadata `json:"metadata"`
}

func NewListMessagesHTTPRequest(c *gin.Context) *ListMessagesRequest {
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit == 0 || limit > 1000 {
		limit = 100
	}

	offset, _ := strconv.Atoi(c.Query("offset"))
	status := c.QueryArray("status")
	if len(status) == 0 {
		status = nil
	}

	return &ListMessagesRequest{
		Limit:  limit,
		Offset: offset,
		Status: status,
	}
}
