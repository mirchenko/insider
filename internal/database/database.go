package database

import (
	"insider/config"
	"insider/pkg/logger"

	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func NewDatabase(cfg *config.Config, log *logger.Logger, shutdowner fx.Shutdowner) (*gorm.DB, error) {
	dbLogger := gormlogger.New(
		log,
		gormlogger.Config{},
	)

	if cfg.DatabaseConfig.Debug {
		dbLogger = dbLogger.LogMode(gormlogger.Info)
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseConfig.URL), &gorm.Config{
		Logger: dbLogger,
	})

	if err != nil {
		log.Error().Err(err).Msg("failed to ping database")
		_ = shutdowner.Shutdown()
		return nil, err
	}
	return db, nil
}
