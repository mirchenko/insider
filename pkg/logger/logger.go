package logger

import (
	"os"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

type Logger struct {
	zerolog.Logger
}

func New() *Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &Logger{Logger: logger}
}

var Module = fx.Provide(New)
