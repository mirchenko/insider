package http

import (
	"context"
	"encoding/json"
	"insider/config"
	"insider/internal/model"
	"insider/internal/repository"
	"insider/internal/sender"
	"insider/internal/services"
	"insider/mocks"
	mocksrepository "insider/mocks/repository"
	"insider/pkg/logger"
	http2 "net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/require"
	"github.com/stretchr/testify/mock"
	"go.uber.org/fx"
)

var messagesTestMockSender *mocks.MockSender
var messagesTestMockMessagesRepo *mocksrepository.MockMessageRepository

func mockAppWithMessagesService() (*fx.App, *MessagesHandler) {
	messagesTestMockSender = new(mocks.MockSender)
	messagesTestMockMessagesRepo = new(mocksrepository.MockMessageRepository)

	var handler *MessagesHandler

	app := fx.New(
		fx.Supply(fx.Annotate(messagesTestMockMessagesRepo, fx.As(new(repository.MessageRepository)))),
		fx.Supply(fx.Annotate(messagesTestMockSender, fx.As(new(sender.Sender)))),
		services.Module,
		fx.Provide(NewMessagesHandler),
		config.Module,
		logger.Module,
		fx.Populate(&handler),
	)

	return app, handler
}

func TestMessagesHandler_ListMessages_Success(t *testing.T) {
	app, handler := mockAppWithMessagesService()

	var messages []model.Message
	var total int64
	messagesTestMockMessagesRepo.EXPECT().Fetch(mock.Anything, mock.Anything).Return(messages, nil)
	messagesTestMockMessagesRepo.EXPECT().Count(mock.Anything, mock.Anything).Return(&total, nil)

	require.NoError(t, app.Start(context.Background()))
	defer app.Stop(context.Background())

	r := gin.New()
	r.GET("/messages", handler.ListMessages)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(
		http2.MethodGet,
		"/messages",
		nil,
	)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	require.Equal(t, http2.StatusOK, w.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	_, hasData := body["data"]
	meta, hasMeta := body["metadata"]
	require.True(t, hasData, "response should have 'data'")
	require.True(t, hasMeta, "response should have 'metadata'")
	m, ok := meta.(map[string]any)
	require.True(t, ok, "metadata should be an object")
	_, hasLimit := m["limit"]
	_, hasOffset := m["offset"]
	_, hasCount := m["count"]
	require.True(t, hasLimit, "metadata.limit missing")
	require.True(t, hasOffset, "metadata.offset missing")
	require.True(t, hasCount, "metadata.count missing")
}

func TestMessagesHandler_ListMessages_InvalidStatus(t *testing.T) {
	app, handler := mockAppWithMessagesService()
	require.NoError(t, app.Start(context.Background()))
	defer app.Stop(context.Background())

	r := gin.New()
	r.GET("/messages", handler.ListMessages)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(
		http2.MethodGet,
		"/messages?status=sent,sent1",
		nil,
	)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	require.Equal(t, http2.StatusBadRequest, w.Code)
}
