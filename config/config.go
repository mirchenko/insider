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
	CacheConfig
}
type DatabaseConfig struct {
	PostgresConfig
}

type PostgresConfig struct {
	URL   string `yaml:"url"`
	Debug bool   `yaml:"debug"`
}

type RedisConfig struct {
	URL string `yaml:"url"`
}

type CacheConfig struct {
	RedisConfig
}

type WebhookProviderConfig struct {
	BaseURL             string `yaml:"base_url"`
	TenantID            string `yaml:"tenant_id"`
	RetriesCount        int    `yaml:"retries_count"`
	RetryTimeoutSeconds int    `yaml:"retry_timeout_seconds"`
	AuthKey             string `yaml:"auth_key"`
	Debug               bool   `yaml:"debug"`
}

type SenderConfig struct {
	IterDurationSeconds int `yaml:"iter_duration_seconds"`
	IterBufferSize      int `yaml:"iter_buffer_size"`
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
			PostgresConfig{
				URL:   viper.GetString("database.postgres.url"),
				Debug: viper.GetBool("database.postgres.debug"),
			},
		},
		WebhookProviderConfig: WebhookProviderConfig{
			BaseURL:             viper.GetString("providers.webhook.base_url"),
			TenantID:            viper.GetString("providers.webhook.tenant_id"),
			RetriesCount:        viper.GetInt("providers.webhook.retries_count"),
			RetryTimeoutSeconds: viper.GetInt("providers.webhook.retry_timeout_seconds"),
			AuthKey:             viper.GetString("providers.webhook.auth_key"),
			Debug:               viper.GetBool("providers.webhook.debug"),
		},
		SenderConfig: SenderConfig{
			IterDurationSeconds: viper.GetInt("sender.iter_duration_seconds"),
			IterBufferSize:      viper.GetInt("sender.iter_buffer_size"),
		},
		CacheConfig: CacheConfig{
			RedisConfig{
				URL: viper.GetString("cache.redis.url"),
			},
		},
	}
}

var Module = fx.Provide(LoadConfig)
