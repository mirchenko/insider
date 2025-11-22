package config

import (
	"insider/pkg/logger"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type Config struct {
	Port string
	DatabaseConfig
	WebhookProviderConfig
	SenderConfig
}
type DatabaseConfig struct {
	URL string
}

type WebhookProviderConfig struct {
	BaseURL             string `yaml:"base_url"`
	TenantID            string `yaml:"tenant_id"`
	RetriesCount        int    `yaml:"retries_count"`
	RetryTimeoutSeconds int    `yaml:"retry_timeout_seconds"`
	AuthKey             string `yaml:"auth_key"`
}

type SenderConfig struct {
	CycleDurationSeconds int `yaml:"cycle_duration_seconds"`
}

func LoadConfig(log *logger.Logger, shutdowner fx.Shutdowner) *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Error().Err(err).Msg("failed to read config file")
		if err := shutdowner.Shutdown(); err != nil {
			log.Error().Err(err).Msg("failed to shutdown application")
		}
	}

	return &Config{
		Port: viper.GetString("port"),
		DatabaseConfig: DatabaseConfig{
			URL: viper.GetString("database.url"),
		},
		WebhookProviderConfig: WebhookProviderConfig{
			BaseURL:             viper.GetString("providers.webhook.base_url"),
			TenantID:            viper.GetString("providers.webhook.tenant_id"),
			RetriesCount:        viper.GetInt("providers.webhook.retries_count"),
			RetryTimeoutSeconds: viper.GetInt("providers.webhook.retry_timeout_seconds"),
			AuthKey:             viper.GetString("providers.webhook.auth_key"),
		},
		SenderConfig: SenderConfig{
			CycleDurationSeconds: viper.GetInt("sender.cycle_duration_seconds"),
		},
	}
}

var Module = fx.Provide(LoadConfig)
