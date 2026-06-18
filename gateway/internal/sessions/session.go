package sessions

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// the store
type Store interface {
	UpgradeTokenVersion(ctx context.Context, userID string) error
	GetTokenVersion(ctx context.Context, userID string) (int64, error)
	ClaimRefreshToken(ctx context.Context, tokenID string, ttl time.Duration) (bool, error)
}

//redis store
type redisStore struct {
	client *redis.Client
}

//returns a new store
func NewStore(c *redis.Client) Store {
	return &redisStore{
		client:c,
	}
}

//upgrade token version 
func( r *redisStore)UpgradeTokenVersion(ctx context.Context,userID string)error{
	key:=fmt.Sprintf("auth:user:%s:version",userID)
	err:=r.client.Incr(ctx,key).Err()
	if err!=nil{
		return fmt.Errorf("redis increment version error: %w",err)
	}
	return nil
}

//get token version
func (r *redisStore) GetTokenVersion(ctx context.Context, userID string) (int64, error) {
	key := fmt.Sprintf("auth:user:%s:version", userID)
	version, err := r.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("getting token version: %w", err)
	}
	return version, nil
}

func (r *redisStore) ClaimRefreshToken(ctx context.Context, tokenID string, ttl time.Duration) (bool, error) {
	key := fmt.Sprintf("claim:%s", tokenID)

	set, err := r.client.SetNX(ctx, key, "claimed", ttl).Result()
	if err != nil {
		return false, fmt.Errorf("redis setnx error: %w", err)
	}
	return !set, nil
}





