package services

import (
	"context"
	"errors"
	"insider/internal/sender"
	"insider/mocks"
	"testing"

	"github.com/go-openapi/testify/v2/require"
	"go.uber.org/fx"
)

var mockSender *mocks.MockSender

func mockAppWithSenderService() (*fx.App, SenderService) {
	mockSender = new(mocks.MockSender)

	var handler SenderService

	app := fx.New(
		fx.Supply(fx.Annotate(mockSender, fx.As(new(sender.Sender)))),
		fx.Provide(fx.Annotate(NewSenderService, fx.As(new(SenderService)))),
		fx.Populate(&handler),
	)

	return app, handler
}

func TestSenderHandler_Toggle_Start_Success(t *testing.T) {
	app, handler := mockAppWithSenderService()

	mockSender.EXPECT().Start().Return(nil)
	mockSender.EXPECT().IsStarted().Return(false)

	require.NoError(t, app.Start(context.Background()))
	defer app.Stop(context.Background())

	res, err := handler.Toggle()

	require.NoError(t, err)

	require.Equal(t, res, "started", "expected started, got", res)
}

func TestSenderHandler_Toggle_Stop_Success(t *testing.T) {
	app, handler := mockAppWithSenderService()

	mockSender.EXPECT().Stop().Return()
	mockSender.EXPECT().IsStarted().Return(true)

	require.NoError(t, app.Start(context.Background()))
	defer app.Stop(context.Background())

	res, err := handler.Toggle()

	require.NoError(t, err)

	require.Equal(t, res, "stopped", "expected stopped, got", res)
}

func TestSenderHandler_Toggle_Start_Error(t *testing.T) {
	app, handler := mockAppWithSenderService()

	mockSender.EXPECT().Start().Return(errors.New("error"))
	mockSender.EXPECT().IsStarted().Return(false)

	require.NoError(t, app.Start(context.Background()))
	defer app.Stop(context.Background())

	res, err := handler.Toggle()

	require.Error(t, err)

	require.Equal(t, res, "error", "expected error, got", res)
}
