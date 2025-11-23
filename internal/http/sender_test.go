package http

import (
	"context"
	"encoding/json"
	"insider/config"
	"insider/internal/sender"
	"insider/internal/services"
	"insider/mocks"
	"insider/pkg/logger"
	http2 "net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/require"
	"go.uber.org/fx"
)

var mockSender *mocks.MockSender

func mockAppWithSenderService() (*fx.App, *SenderHandler) {
	mockSender = new(mocks.MockSender)
	mockSender.EXPECT().IsStarted().Return(false)
	mockSender.EXPECT().Start().Return(nil)

	var handler *SenderHandler

	app := fx.New(
		fx.Supply(fx.Annotate(mockSender, fx.As(new(sender.Sender)))),
		services.Module,
		fx.Provide(NewSenderHandler),
		config.Module,
		logger.Module,
		fx.Populate(&handler),
	)

	return app, handler
}

func TestSenderHandler_Toggle_Success(t *testing.T) {
	app, handler := mockAppWithSenderService()
	require.NoError(t, app.Start(context.Background()))
	defer app.Stop(context.Background())

	r := gin.New()
	r.POST("/sender/toggle", handler.ToggleSender)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(
		http2.MethodPost,
		"/sender/toggle",
		nil,
	)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	require.Equal(t, http2.StatusOK, w.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	_, hasStatus := body["status"]
	require.True(t, hasStatus, "response should have 'status'")
	require.Equal(t, body["status"], "started")
}
