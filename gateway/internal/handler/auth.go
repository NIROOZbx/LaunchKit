package handler

import (
	"errors"
	"fmt"

	"github.com/Launchkit-org/LaunchKit/gateway/consts"
	"github.com/Launchkit-org/LaunchKit/gateway/internal/domain"
	"github.com/Launchkit-org/LaunchKit/gateway/internal/models"
	"github.com/Launchkit-org/LaunchKit/gateway/jwt"
	"github.com/Launchkit-org/LaunchKit/shared/apperrors"
	"github.com/Launchkit-org/LaunchKit/shared/config"
	response "github.com/Launchkit-org/LaunchKit/shared/responses"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	authService domain.AuthService
	logger      zerolog.Logger
	jwtCfg      *config.JwtConfig
}

func NewAuthHandler(authService domain.AuthService, logger zerolog.Logger, jwtCfg *config.JwtConfig) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
		jwtCfg:      jwtCfg,
	}
}

func (h *AuthHandler) GetNonce(c fiber.Ctx) error {
	var query models.NonceRequest
	if err := c.Bind().Query(&query); err != nil {
		return response.BadRequest(c, "invalid query parameters", nil)
	}

	nonce, siweMsg, err := h.authService.GetNonce(c.Context(), query.WalletAddress)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get nonce")
		return response.InternalServerError(c)
	}

	return response.OK(c, "Nonce fetched successfully", models.NonceResponse{
		Nonce:   nonce,
		Message: siweMsg,
	})
}

func (h *AuthHandler) Verify(c fiber.Ctx) error {
	payload := new(models.VerifyRequest)
	if err := c.Bind().Body(payload); err != nil {
		return response.BadRequest(c, "invalid request body", nil)
	}

	userResp, tokenPayload, err := h.authService.Verify(c.Context(), payload.Message, payload.Signature, payload.UserType)
	if err != nil {
		h.logger.Error().Err(err).Msg("verification failed")
		if errors.Is(err, apperrors.ErrUnauthorized) {
			return response.Unauthorized(c, "nonce expired or invalid")
		}
		return response.Unauthorized(c, fmt.Sprintf("signature verification failed: %s", err.Error()))
	}

	// Generate JWT pair
	jwtConfig := jwt.FromSharedConfig(h.jwtCfg)
	pair, err := jwt.GenerateTokenPair(*tokenPayload, jwtConfig)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", tokenPayload.UserID).Msg("failed to generate token pair")
		return response.InternalServerError(c)
	}

	isProd := h.jwtCfg.Environment == "production"
	jwt.SetTokenCookies(c, pair, h.jwtCfg.AccessExpiryMinutes, h.jwtCfg.RefreshExpiryHours, isProd)

	return response.OK(c, "Signature verified", userResp)
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	userID, ok := c.Locals(consts.UID).(string)
	if !ok || userID == "" {
		jwt.ClearTokenCookies(c)
		return response.OK(c, "Logged out successfully", nil)
	}

	refreshToken := c.Cookies("refresh_token")

	err := h.authService.Logout(c.Context(), userID, refreshToken)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID).Msg("failed to process logout in service")
	}

	jwt.ClearTokenCookies(c)
	return response.OK(c, "Logged out successfully", nil)
}