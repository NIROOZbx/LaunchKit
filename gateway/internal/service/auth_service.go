package service

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/Launchkit-org/LaunchKit/gateway/internal/constants"
	"github.com/Launchkit-org/LaunchKit/gateway/internal/domain"
	"github.com/Launchkit-org/LaunchKit/gateway/internal/dtos"
	"github.com/Launchkit-org/LaunchKit/gateway/internal/sessions"
	"github.com/Launchkit-org/LaunchKit/gateway/jwt"
	"github.com/Launchkit-org/LaunchKit/shared/apperrors"
	"github.com/Launchkit-org/LaunchKit/shared/config"
	"github.com/spruceid/siwe-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authService struct {
	authStore domain.AuthStore
	sessions  sessions.Store
	client    domain.CoreClient
	jwtCfg    *config.JwtConfig
}

func NewAuthService(authStore domain.AuthStore, sessions sessions.Store, client domain.CoreClient, jwtCfg *config.JwtConfig) domain.AuthService {
	return &authService{
		authStore: authStore,
		sessions:  sessions,
		client:    client,
		jwtCfg:    jwtCfg,
	}
}

func (s *authService) GetNonce(ctx context.Context, walletAddress string) (string, string, error) {
	nonce := siwe.GenerateNonce()
	err := s.authStore.SaveNonce(ctx, nonce, constants.NonceTTL)
	if err != nil {
		return "", "", fmt.Errorf("authService.GetNonce: %w", err)
	}

	domainStr := extractDomain(s.jwtCfg.FrontendURL)

	siweMsg, err := siwe.InitMessage(
		domainStr,
		walletAddress,
		s.jwtCfg.FrontendURL,
		nonce,
		map[string]interface{}{
			"statement": "I accept the LaunchKit Terms of Service",
			"chainId":   1,
			"issuedAt":  time.Now().UTC().Format(time.RFC3339),
		},
	)
	if err != nil {
		return "", "", fmt.Errorf("authService.GetNonce: init SIWE message: %w", err)
	}

	return nonce, siweMsg.String(), nil
}

func (s *authService) Verify(ctx context.Context, message, signature string, userType string) (*dtos.UserResponse, *jwt.TokenPayload, error) {
	siweMessage, err := siwe.ParseMessage(message)
	if err != nil {
		return nil, nil, fmt.Errorf("authService.Verify: parse message: %w", err)
	}

	expectedDomain := extractDomain(s.jwtCfg.FrontendURL)
	nonceVal := siweMessage.GetNonce()

	exists, err := s.authStore.ConsumeNonce(ctx, nonceVal)
	if err != nil {
		return nil, nil, fmt.Errorf("authService.Verify: consume nonce error: %w", err)
	}
	if !exists {
		return nil, nil, apperrors.ErrUnauthorized
	}

	_, err = siweMessage.Verify(signature, &expectedDomain, &nonceVal, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("authService.Verify: signature verification failed: %w", err)
	}

	walletAddress := siweMessage.GetAddress().String()

	// Get or Create user via gRPC client
	user, err := s.client.GetUserByWallet(ctx, walletAddress)
	if err != nil && status.Code(err) != codes.NotFound {
		return nil, nil, fmt.Errorf("authService.Verify: get user error: %w", err)
	}

	if user == nil || (err != nil && status.Code(err) == codes.NotFound) {
		if userType == "" {
			userType = "b2c" // Default to B2C
		}
		user, err = s.client.CreateUser(ctx, walletAddress, userType)
		if err != nil {
			return nil, nil, fmt.Errorf("authService.Verify: create user error: %w", err)
		}
	}

	newVer, err := s.sessions.GetTokenVersion(ctx, user.Id)
	if err != nil {
		return nil, nil, fmt.Errorf("authService.Verify: get token version error: %w", err)
	}

	isOnboarded := domain.IsUserOnboarded(user)

	userResponse := &dtos.UserResponse{
		ID:            user.Id,
		WalletAddress: user.WalletAddress,
		EnsName:       user.EnsName,
		DisplayName:   user.DisplayName,
		AvatarURL:     user.AvatarUrl,
		Role:          user.UserType,
		IsOnboarded:   isOnboarded,
	}

	tokenPayload := &jwt.TokenPayload{
		UserID:        user.Id,
		WalletAddress: user.WalletAddress,
		Role:          user.UserType,
		ProjectID:     user.ProjectId,
		ProjectRole:   user.ProjectRole,
		IsOnboarded:   isOnboarded,
		Version:       newVer,
	}

	return userResponse, tokenPayload, nil
}

func (s *authService) Logout(ctx context.Context, userID, refreshToken string) error {
	err := s.sessions.UpgradeTokenVersion(ctx, userID)
	if err != nil {
		return fmt.Errorf("authService.Logout: upgrade token version error: %w", err)
	}

	if refreshToken != "" {
		claims, err := jwt.ParseRefreshToken(refreshToken, []byte(s.jwtCfg.RefreshTokenSecret))
		if err == nil && claims != nil {
			remaining := time.Until(claims.ExpiresAt.Time)
			if remaining > 0 {
				_, _ = s.sessions.ClaimRefreshToken(ctx, claims.TokenID, remaining)
			}
		}
	}

	return nil
}

func extractDomain(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Host == "" {
		return "localhost"
	}
	return parsed.Host
}
