package services

import (
	"context"
	"insider/internal/model"
	"insider/internal/repository"
	mocksrepository "insider/mocks/repository"
	"testing"

	"github.com/go-openapi/testify/v2/require"
	"github.com/stretchr/testify/mock"
	"go.uber.org/fx"
)

var mockMessagesRepo *mocksrepository.MockMessageRepository

func mockAppWithMessagesService() (*fx.App, MessagesService) {
	mockMessagesRepo = new(mocksrepository.MockMessageRepository)

	var handler MessagesService

	app := fx.New(
		fx.Supply(fx.Annotate(mockMessagesRepo, fx.As(new(repository.MessageRepository)))),
		fx.Provide(fx.Annotate(NewMessagesService, fx.As(new(MessagesService)))),
		fx.Populate(&handler),
	)

	return app, handler
}

func TestMessagesHandler_List_Success(t *testing.T) {
	app, handler := mockAppWithMessagesService()

	messages := []model.Message{
		{
			ID: 1,
		},
	}
	total := int64(100)
	mockMessagesRepo.EXPECT().Fetch(mock.Anything, mock.Anything).Return(messages, nil)
	mockMessagesRepo.EXPECT().Count(mock.Anything, mock.Anything).Return(&total, nil)

	require.NoError(t, app.Start(context.Background()))
	defer app.Stop(context.Background())

	res, err := handler.List(context.Background(), &model.ListMessagesRequest{})

	require.NoError(t, err)

	require.Equal(t, res.Count, total, "expected 100, got", total)
	require.Equal(t, res.Data[0].ID, int64(1), "expected 1, got", res.Data[0].ID)
}
