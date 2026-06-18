package store

import (
	"context"
	"time"

	"github.com/Launchkit-org/LaunchKit/gateway/internal/constants"
	"github.com/Launchkit-org/LaunchKit/shared/cache"
)

type RedisAuthStore struct {
	cache cache.Cache
}

func NewRedisAuthStore(cache cache.Cache) *RedisAuthStore {
	return &RedisAuthStore{
		cache: cache,
	}
}

func (s *RedisAuthStore) SaveNonce(ctx context.Context, nonce string, ttl time.Duration) error {
	key := constants.RedisKeyPrefixNonce + nonce
	return s.cache.Set(ctx, key, "true", ttl)
}

func (s *RedisAuthStore) ConsumeNonce(ctx context.Context, nonce string) (bool, error) {
	key := constants.RedisKeyPrefixNonce + nonce
	return s.cache.Delete(ctx, key)
}
