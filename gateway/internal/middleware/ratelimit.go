package middleware

import (
	_ "embed"
	"fmt"
	"time"

	response "github.com/Launchkit-org/LaunchKit/shared/responses"
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
)

//go:embed ratelimit.lua
var rateLimitScriptSrc string

var rateLimitScript = redis.NewScript(rateLimitScriptSrc)

// IdentifierFunc extracts the key to rate-limit on (IP, user id, etc.)
// from the request. Return an error to reject the request outright
// (e.g. protected route hit with no identity in context).
type IdentifierFunc func(c fiber.Ctx) (string, error)

type RateLimiterConfig struct {
	Redis      *redis.Client
	Limit      int
	Window     time.Duration
	KeyPrefix  string
	Identifier IdentifierFunc
}

func IdentifierFromIP(c fiber.Ctx) (string, error) {
	ip := c.IP()
	if ip == "" {
		return "", fmt.Errorf("no client ip")
	}
	return ip, nil
}

// IdentifierFromUser pulls the authenticated identity set by the JWT
// middleware.
func IdentifierFromUser(c fiber.Ctx) (string, error) {
	v := c.Locals("user_id")
	id, ok := v.(string)
	if !ok || id == "" {
		return "", fmt.Errorf("no authenticated user in context")
	}
	return id, nil
}

func NewRateLimiter(cfg RateLimiterConfig) fiber.Handler {
	windowSeconds := int64(cfg.Window.Seconds())

	return func(c fiber.Ctx) error {
		id, err := cfg.Identifier(c)
		if err != nil {
			return response.Unauthorized(c, "unable to identify client for rate limiting")
		}

		now := time.Now().Unix()
		currentWindow := now / windowSeconds
		previousWindow := currentWindow - 1
		elapsed := now % windowSeconds

		currKey := fmt.Sprintf("%s:%s:%d", cfg.KeyPrefix, id, currentWindow)
		prevKey := fmt.Sprintf("%s:%s:%d", cfg.KeyPrefix, id, previousWindow)

		weight := float64(windowSeconds-elapsed) / float64(windowSeconds)

		res, err := rateLimitScript.Run(
			c.Context(),
			cfg.Redis,
			[]string{currKey, prevKey},
			weight, cfg.Limit, windowSeconds*2,
		).Result()

		if err != nil {
			// Fail-open: don't take the whole API down on a Redis blip.
			// TODO: log this via a.Logger once wired through.
			return c.Next()
		}

		allowed, _ := res.(int64)
		if allowed == 0 {
			c.Set("Retry-After", fmt.Sprintf("%d", windowSeconds-elapsed))
			return response.TooManyRequests(c, "rate limit exceeded, slow down")
		}

		return c.Next()
	}
}