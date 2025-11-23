package cache

import (
	"context"
	"insider/config"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(cfg *config.Config, shutdowner fx.Shutdowner) (*Redis, error) {
	opts, err := redis.ParseURL(cfg.CacheConfig.URL)
	if err != nil {
		log.Error().Err(err).Msg("failed to init redis")
		_ = shutdowner.Shutdown()
		return nil, err
	}

	return &Redis{client: redis.NewClient(opts)}, nil
}

func (r *Redis) Set(ctx context.Context, key, value string) error {
	return r.client.Set(ctx, key, value, 0).Err()
}
