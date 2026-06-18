package constants

import "time"

// Redis key prefixes
const (
	RedisKeyPrefixNonce     = "auth:nonce:"
	RedisKeyPrefixRateLimit = "ratelimit:"

	RateLimitKeyPrefixAuth   = "ratelimit:auth"
	RateLimitKeyPrefixPublic = "ratelimit:public"
	RateLimitKeyPrefixAPI    = "ratelimit:api"
)

// Cache TTLs
const (
	NonceTTL = 5 * time.Minute
)

// Cookie names
const (
	CookieAccessToken  = "access_token"
	CookieRefreshToken = "refresh_token"
)

// Fiber context locals keys
const (
	UserID        = "user_id"
	WalletAddress = "wallet_address"
	Role          = "role"
	IsOnboarded   = "is_onboarded"
)

// Rate limit windows (seconds)
const (
	RateLimitWindowSec int64 = 60
)
