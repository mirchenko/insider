package cache

import "context"

type Cache interface {
	Set(ctx context.Context, key, value string) error
}
