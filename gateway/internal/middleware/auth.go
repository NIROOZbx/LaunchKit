package middleware

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/Launchkit-org/LaunchKit/gateway/consts"
	"github.com/Launchkit-org/LaunchKit/gateway/internal/domain"
	"github.com/Launchkit-org/LaunchKit/gateway/internal/sessions"
	"github.com/Launchkit-org/LaunchKit/gateway/jwt"
	"github.com/Launchkit-org/LaunchKit/shared/apperrors"
	"github.com/Launchkit-org/LaunchKit/shared/config"
	response "github.com/Launchkit-org/LaunchKit/shared/responses"
	"github.com/gofiber/fiber/v3"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

type authMiddleware struct {
	cfg   *config.JwtConfig
	log   zerolog.Logger
	repo  domain.CoreClient
	store sessions.Store
}

type AuthMiddleware interface {
	Auth(c fiber.Ctx) error
	OnboardingAuth(c fiber.Ctx) error
	RequireRole(allowedRoles ...string) fiber.Handler
}

func NewMiddleware(store sessions.Store, cfg *config.JwtConfig, log zerolog.Logger, repo domain.CoreClient) AuthMiddleware {
	return &authMiddleware{
		store: store,
		cfg:   cfg,
		log:   log,
		repo:  repo,
	}
}

const graceperiod = time.Second * 5

func (a *authMiddleware) Auth(c fiber.Ctx) error {
	claims, err := a.validate(c)
	if err != nil {
		a.log.Warn().Err(err).Msg("authentication failed: invalid or missing token")
		return response.Unauthorized(c, apperrors.ErrForbidden.Error())
	}
	if claims.Role == "b2b" && !claims.IsOnboarded {
		a.log.Info().Str("userID", claims.UserID).Msg("workspace setup required for B2B user")
		return response.Forbidden(c, "workspace setup required", nil)
	}
	c.Locals(consts.UID, claims.UserID)
	c.Locals(consts.PRID, claims.ProjectID)
	c.Locals(consts.Role, claims.Role)

	a.log.Debug().
		Str("userID", claims.UserID).
		Str("projectID", claims.ProjectID).
		Str("role", claims.Role).
		Msg("user authenticated successfully")

	return c.Next()
}

// onboarding auth check
func (a *authMiddleware) OnboardingAuth(c fiber.Ctx) error {

	claims, err := a.validate(c)
	if err != nil {
		fmt.Println("err", err)
		return response.Unauthorized(c, apperrors.ErrUnauthorized.Error())
	}
	if claims.IsOnboarded {
		a.log.Info().Str("userID", claims.UserID).Msg("onboarding auth attempt for already onboarded user")
		return response.Forbidden(c, "already onboarded", nil)
	}
	c.Locals(consts.UID, claims.UserID)
	c.Locals(consts.Role, claims.Role)

	return c.Next()

}

//validate access

func (a *authMiddleware) validate(c fiber.Ctx) (*jwt.AccessClaims, error) {

	accessToken := c.Cookies("access_token")
	if accessToken == "" {
		a.log.Debug().Msg("access token cookie missing, attempting silent refresh")
		return a.silentRefresh(c)
	}

	claims, err := jwt.ParseAccessToken(accessToken, []byte(a.cfg.AccessTokenSecret))
	if errors.Is(err, gojwt.ErrTokenExpired) {
		if claims == nil {
			a.log.Debug().Msg("access token expired no attempt")
			return nil, apperrors.ErrUnauthorized
		}
		a.log.Debug().Msg("access token expired, attempting silent refresh")
		return a.silentRefresh(c)
	}

	if err != nil {
		a.log.Debug().Err(err).Msg("access token invalid, attempting recovery via refresh")
		return a.silentRefresh(c)
	}
	version, verErr := a.store.GetTokenVersion(c.Context(), claims.UserID)
	if verErr != nil {
		return claims, nil
	}
	if claims.Version < version {
		a.log.Debug().Msg("token version mismatch")
		jwt.ClearTokenCookies(c)
		return nil, apperrors.ErrUnauthorized
	}
	return claims, nil
}

//silent refresh token issue

func (a *authMiddleware) silentRefresh(c fiber.Ctx) (*jwt.AccessClaims, error) {

	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		jwt.ClearTokenCookies(c)
		return nil, apperrors.ErrUnauthorized
	}
	refreshClaims, err := jwt.ParseRefreshToken(refreshToken, []byte(a.cfg.RefreshTokenSecret))
	if err != nil {
		a.log.Warn().Err(err).Msg("failed to parse refresh token during silent refresh")
		jwt.ClearTokenCookies(c)
		return nil, apperrors.ErrUnauthorized
	}

	userID := refreshClaims.UserID

	User, err := a.repo.GetUser(c.Context(), userID)
	if err != nil {
		a.log.Error().Err(err).Str("userID", refreshClaims.UserID).Msg("failed to fetch user  from DB")
		return nil, apperrors.ErrUnauthorized
	}

	alreadyClaimed, err := a.store.ClaimRefreshToken(c.Context(), refreshClaims.TokenID, 30*time.Second)

	if err != nil {
		a.log.Error().Err(err).Str("userID", refreshClaims.UserID).Msg("redis down during token claim")
		return nil, apperrors.ErrInternal
	}

	if alreadyClaimed {
		a.log.Debug().Str("tokenID", refreshClaims.TokenID).Msg("token already claimed by another request")
		return nil, apperrors.ErrRefreshLockHeld
	}

	newVer, err := a.store.GetTokenVersion(c.Context(), refreshClaims.UserID)
	if err != nil {
		a.log.Error().Err(err).Str("userID", refreshClaims.UserID).Msg("failed to fetch token version during silent refresh")
	}
	isOnboarded := domain.IsUserOnboarded(User)

	payload := &jwt.TokenPayload{
		UserID:        refreshClaims.UserID,
		WalletAddress: User.WalletAddress,
		Role:          User.UserType,
		ProjectID:     User.ProjectId,
		ProjectRole:   User.ProjectRole,
		IsOnboarded:   isOnboarded,
		Version:       newVer,
	}
	jwtConfig := a.cfg.ToJWTConfig()

	pair, err := jwt.GenerateTokenPair(*payload, jwtConfig)
	if err != nil {
		a.log.Error().Err(err).Str("userID", refreshClaims.UserID).Msg("failed to generate token pair during silent refresh")
		return nil, err
	}
	isProd := a.cfg.Environment == "production"

	jwt.SetTokenCookies(c, pair, a.cfg.AccessExpiryMinutes, a.cfg.RefreshExpiryHours, isProd)

	a.log.Info().
		Str("userID", refreshClaims.UserID).
		Str("ProjectID", User.ProjectId).
		Msg("tokens issued successfully during silent refresh")

	return jwt.ParseAccessToken(pair.AccessToken, []byte(a.cfg.AccessTokenSecret))
}

func (a *authMiddleware) RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {

		role, ok := c.Locals(consts.Role).(string)
		if !ok {
			a.log.Warn().Msg("role middleware: role not found in context")
			return response.Forbidden(c, "role not found", nil)
		}
		if slices.Contains(allowedRoles, role) {
			return c.Next()
		}
		a.log.Warn().
			Str("userRole", role).
			Interface("allowedRoles", allowedRoles).
			Msg("role middleware: insufficient permissions")
		return response.Forbidden(c, "insufficient permissions", nil)
	}
}
