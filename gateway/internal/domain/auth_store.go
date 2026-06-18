package domain

import (
	"context"
	"time"
)

type AuthStore interface {
	SaveNonce(ctx context.Context, nonce string, ttl time.Duration) error
	ConsumeNonce(ctx context.Context, nonce string) (bool, error)
}
