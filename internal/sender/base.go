package sender

import (
	"context"
	"errors"
	"insider/config"
	"insider/internal/model"
	"insider/internal/provider"
	"insider/internal/repository"
	"insider/pkg/logger"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type Sender struct {
	provider provider.Provider
	logger   *logger.Logger
	repo     *repository.MessagesRepository
	cfg      *config.Config
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewSender(provider *provider.WebhookProvider, log *logger.Logger, cfg *config.Config, repo *repository.MessagesRepository) *Sender {
	return &Sender{
		provider: provider,
		logger:   log,
		repo:     repo,
		cfg:      cfg,
	}
}

func (s *Sender) Start() error {
	if s.ctx != nil && s.ctx.Err() == nil {
		return errors.New("sender already started in this instance")
	}
	s.logger.Info().Msg("starting sender")
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.run()

	return nil
}

func (s *Sender) Stop() {
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
}

func (s *Sender) run() {
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				s.logger.Info().Msg("stopping sender")
				return
			default:
				s.logger.Info().Msg("sender iteration")
				if err := s.handle(ctx); err != nil {
					s.logger.Error().Err(err).Msg("error while handling message")
				}
				time.Sleep(time.Duration(s.cfg.CycleDurationSeconds) * time.Second)
			}

		}
	}(s.ctx)
}

func (s *Sender) handle(ctx context.Context) error {
	res, err := s.repo.Fetch(ctx, &model.ListMessagesRequest{
		Limit:  2,
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

			res, err := s.provider.Send(&v)
			if err != nil {
				v.Status = model.MessageStatusFailed
				e := err.Error()
				v.Reason = &e
			} else {
				v.Status = model.MessageStatusSent
				v.ExternalMessageID = res.MessageID
			}

			if _, err := s.repo.Update(ctx, &v); err != nil {
				s.logger.Error().Err(err)
			}
			// TODO: write to redis
		}()
	}

	wg.Wait()
	return nil
}
