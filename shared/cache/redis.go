package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Lauchkit/LaunchKit/shared/config"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg *config.RedisConfig) (*redis.Client,error){

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
		Password: cfg.Password,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}
	log.Println("Redis connected ✅")

	return rdb, nil
}