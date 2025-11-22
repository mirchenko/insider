package repository

import (
	"context"
	"insider/internal/model"
	"insider/pkg/logger"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type MessagesRepository struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewMessagesRepository(db *gorm.DB, log *logger.Logger) *MessagesRepository {
	return &MessagesRepository{db: db, log: log}
}

func (m *MessagesRepository) Fetch(ctx context.Context, r *model.ListMessagesRequest) ([]model.Message, error) {
	messages, err := buildListQuery(m.db, r).Limit(r.Limit).Offset(r.Offset).Order("created_at desc").Find(ctx)

	if err != nil {
		m.log.Error().Err(err).Msg("failed to get messages")
		return nil, err
	}
	return messages, nil
}

func (m *MessagesRepository) Count(ctx context.Context, r *model.ListMessagesRequest) (*int64, error) {
	count, err := buildListQuery(m.db, r).Count(ctx, "id")

	if err != nil {
		m.log.Error().Err(err).Msg("failed to count messages")
		return nil, err
	}
	return &count, nil
}

func (m *MessagesRepository) Update(ctx context.Context, msg *model.Message) (*model.Message, error) {
	_, err := gorm.G[model.Message](m.db).Where("id = ?", msg.ID).Updates(ctx, *msg)
	if err != nil {
		m.log.Error().Err(err).Dict("message", zerolog.Dict().Int64("id", msg.ID)).Msg("failed to update message")
		return nil, err
	}

	return msg, nil
}

func buildListQuery(db *gorm.DB, r *model.ListMessagesRequest) gorm.ChainInterface[model.Message] {
	query := gorm.G[model.Message](db).Where("1 = 1")
	if r.Status != nil {
		query = query.Where("status in ?", r.Status)
	}
	return query
}
