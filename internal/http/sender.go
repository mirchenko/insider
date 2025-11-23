package http

import (
	"insider/internal/services"

	"github.com/gin-gonic/gin"
)

type SenderHandler struct {
	svc services.SenderService
}

type ToggleSenderResponse struct {
	Status string `json:"status"`
}

func NewSenderHandler(svc services.SenderService) *SenderHandler {
	return &SenderHandler{svc: svc}
}

// ToggleSender godoc
// @Summary Starts/Stops sender
// @Description starts/stop sender based on sender current status, return started, stopped, error statuses in response
// @Tags Sender
// @Accept json
// @Produce json
// @Success 200 {object} ToggleSenderResponse
// @Example 200 {"status":"started"}
// @Example 200 {"status":"stopped"}
// @Success 409 {object} ToggleSenderResponse
// @Example 409 {"status":"error"}
// @Router /sender/toggle [post]
func (h *SenderHandler) ToggleSender(c *gin.Context) {
	status, err := h.svc.Toggle()
	if err != nil {
		c.JSON(409, ToggleSenderResponse{
			status,
		})

		return
	}

	c.JSON(200, ToggleSenderResponse{
		status,
	})
}
