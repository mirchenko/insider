package http

import (
	"context"
	"insider/config"
	"insider/internal/sender"
	"insider/mocks"
	"insider/pkg/logger"
	http2 "net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/require"
	"go.uber.org/fx"
)

func mockAppWithHealthService() (*fx.App, *HealthHandler) {
	var mockSender *mocks.MockSender
	var handler *HealthHandler

	app := fx.New(
		fx.Supply(fx.Annotate(mockSender, fx.As(new(sender.Sender)))),
		config.Module,
		logger.Module,
		fx.Provide(NewHealthHandler),
		fx.Populate(&handler),
	)

	return app, handler
}

func TestHealth_Success(t *testing.T) {
	app, handler := mockAppWithHealthService()
	require.NoError(t, app.Start(context.Background()))
	defer app.Stop(context.Background())

	r := gin.New()
	r.GET("/health", handler.Health)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(
		http2.MethodGet,
		"/health",
		nil,
	)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	require.Equal(t, http2.StatusOK, w.Code)
}
