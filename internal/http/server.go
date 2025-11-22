package http

import (
	"context"
	"insider/config"
	"insider/internal/http/validators"
	"insider/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func NewServer(routes *Routers, log *logger.Logger, shutdowner fx.Shutdowner) *gin.Engine {
	r := gin.Default()
	if err := validators.RegisterValidators(); err != nil {
		log.Error().Err(err).Msg("failed to register gin http custom validators")
		_ = shutdowner.Shutdown()
	}
	routes.Register(r)
	return r
}

func StartServer(
	lc fx.Lifecycle,
	cfg *config.Config,
	r *gin.Engine,
	shutdowner fx.Shutdowner,
	log *logger.Logger,
) {
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil {
					if err := shutdowner.Shutdown(); err != nil {
						log.Error().Err(err).Msg("failed to shutdown application")
					}
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info().Msg("shutting down HTTP server...")
			if err := srv.Shutdown(ctx); err != nil {
				log.Error().Err(err).Msg("HTTP server shutdown error")
				return err
			}

			log.Info().Msg("HTTP server stopped gracefully")
			return nil
		},
	})
}
