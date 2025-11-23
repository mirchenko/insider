package sender

import (
	"context"
	"errors"
	"insider/config"
	"insider/internal/cache"
	"insider/internal/model"
	"insider/internal/provider"
	"insider/internal/repository"
	"insider/pkg/logger"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type WorkerSender struct {
	provider provider.Provider
	logger   *logger.Logger
	repo     repository.MessageRepository
	cache    cache.Cache
	cfg      *config.Config
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewWorkerSender(provider *provider.WebhookProvider, log *logger.Logger, cfg *config.Config, repo repository.MessageRepository, mc *cache.MessagesCache) *WorkerSender {
	return &WorkerSender{
		provider: provider,
		logger:   log,
		repo:     repo,
		cfg:      cfg,
		cache:    mc,
	}
}

func (s *WorkerSender) IsStarted() bool {
	if s.ctx != nil && s.ctx.Err() == nil {
		return true
	}

	return false
}

func (s *WorkerSender) Start() error {
	if s.IsStarted() {
		return errors.New("sender already started in this instance")
	}
	s.logger.Info().Msg("starting sender")
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.run()

	return nil
}

func (s *WorkerSender) Stop() {
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
}

func (s *WorkerSender) run() {
	go func(ctx context.Context) {
		for {
			s.logger.Info().Msg("sender iteration")
			if err := s.handle(ctx); err != nil {
				s.logger.Error().Err(err).Msg("error while handling messages")
			}

			select {
			case <-ctx.Done():
				s.logger.Info().Msg("stopping sender")
				return
			case <-time.After(time.Duration(s.cfg.IterDurationSeconds) * time.Second):
			}

		}
	}(s.ctx)
}

func (s *WorkerSender) handle(ctx context.Context) error {
	res, err := s.repo.Fetch(ctx, &model.ListMessagesRequest{
		Limit:  s.cfg.IterBufferSize,
		Status: []string{model.MessageStatusPending},
	})
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(res))

	for _, v := range res {
		go func() {
			defer wg.Done()
			s.logger.Info().Dict("message", zerolog.Dict().Int64("id", v.ID)).Msg("processing message")

			sendTime := time.Now()
			res, err := s.provider.Send(&v)
			if err != nil {
				v.Status = model.MessageStatusFailed
				e := err.Error()
				v.Reason = &e
			} else {
				v.Status = model.MessageStatusSent
				v.ExternalMessageID = res.MessageID

				if err := s.cache.Set(ctx, v.ExternalMessageID, sendTime.String()); err != nil {
					s.logger.Error().Err(err)
				}
			}

			if _, err := s.repo.Update(ctx, &v); err != nil {
				s.logger.Error().Err(err)
			}

			s.logger.Info().Dict("message", zerolog.Dict().Int64("id", v.ID)).Msg("message processed")
		}()
	}

	wg.Wait()
	return nil
}
