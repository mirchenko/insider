package http

import (
	"insider/internal/http/validators"
	"insider/internal/model"
	"insider/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

type MessagesHandler struct {
	svc services.MessagesService
}

func NewMessagesHandler(svc services.MessagesService, shutdowner fx.Shutdowner) *MessagesHandler {
	if err := validators.RegisterValidators(); err != nil {
		log.Error().Err(err).Msg("failed to register gin http custom validators")
		_ = shutdowner.Shutdown()
	}

	return &MessagesHandler{svc: svc}
}

// ListMessages godoc
// @Summary list messages
// @Description list messages by limit, offset and status
// @Tags Messages
// @Accept json
// @Produce json
// @Param limit query int false "Limit int" minimum(1) maximum(1000) default(24)
// @Param offset query int false "Offset int" minimum(0) maximum(1000) default(0)
// @Param status query []string false "Status enum(sent, delivered, pending, failed)" enums(sent,delivered,pending,failed) default(sent)
// @Success 200 {array} model.Message
// @Fail 400 {object} gin.H "Validation error"
// @Router /messages [get]
func (h *MessagesHandler) ListMessages(c *gin.Context) {
	var req model.ListMessagesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := h.svc.List(c.Request.Context(), &req)

	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to fetch messages",
		})
		return
	}

	c.JSON(200, gin.H{
		"data":     res.Data,
		"metadata": res.ListResponseMetadata,
	})
}
