package app

import (
	"insider/config"
	"insider/internal/database"
	"insider/internal/handler"
	"insider/internal/http"
	"insider/internal/provider"
	"insider/internal/repository"
	"insider/internal/sender"
	"insider/internal/services"
	"insider/pkg/logger"

	"go.uber.org/fx"
)

var Module = fx.Options(
	logger.Module,
	config.Module,
	database.Module,
	repository.Module,
	services.Module,
	handler.Module,
	http.Module,
	provider.Module,
	sender.Module,
)
