package services

import (
	"context"
	"insider/internal/model"
	"insider/internal/repository"

	"golang.org/x/sync/errgroup"
)

type MessagesService struct {
	repo *repository.MessagesRepository
}

func NewMessagesService(repo *repository.MessagesRepository) *MessagesService {
	return &MessagesService{repo: repo}
}

func (s *MessagesService) List(ctx context.Context, r *model.ListMessagesRequest) (*model.ListMessagesResponse, error) {
	if r.Status == nil {
		r.Status = []string{model.MessageStatusSent}
	}

	var messages []model.Message
	var count *int64

	g := errgroup.Group{}

	g.Go(func() error {
		data, err := s.repo.Fetch(ctx, r)
		if err != nil {
			return err
		}
		messages = data
		return nil
	})

	g.Go(func() error {
		total, err := s.repo.Count(ctx, r)
		if err != nil {
			return err
		}

		count = total
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &model.ListMessagesResponse{
		Data: messages,
		ListResponseMetadata: model.ListResponseMetadata{
			Count:  *count,
			Limit:  r.Limit,
			Offset: r.Offset,
		},
	}, nil
}
