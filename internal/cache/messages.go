package cache

import "context"

const messagesRedisKey = "messages"

type MessagesCache struct {
	redis Cache
}

func NewMessagesCache(redis *Redis) *MessagesCache {
	return &MessagesCache{redis: redis}
}

func (c *MessagesCache) Set(ctx context.Context, key, value string) error {
	return c.redis.Set(ctx, makeRedisKey(messagesRedisKey, key), value)
}
