package handler

import (
	"insider/internal/services"

	"github.com/gin-gonic/gin"
)

type SenderHandler struct {
	svc *services.SenderService
}

func NewSenderHandler(svc *services.SenderService) *SenderHandler {
	return &SenderHandler{svc: svc}
}

// StartSender godoc
// @Summary Starts sender
// @Description starts sender if it isn't running yet
// @Tags Sender
// @Accept json
// @Produce json
// @Success 200 {"ok": true}
// @Fail 409 {object} gin.H "Conflict error"
// @Router /messages [get]
func (h *SenderHandler) StartSender(c *gin.Context) {
	if err := h.svc.Start(); err != nil {
		c.JSON(409, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"ok": true,
	})
}

// StopSender godoc
// @Summary stops sender
// @Description stop sender
// @Tags Sender
// @Accept json
// @Produce json
// @Success 200 {"ok": true}
// @Router /messages [get]
func (h *SenderHandler) StopSender(c *gin.Context) {
	h.svc.Stop()

	c.JSON(200, gin.H{
		"ok": true,
	})
}
