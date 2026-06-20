package domain

import (
	"context"

	"github.com/Launchkit-org/LaunchKit/gateway/internal/models"
	"github.com/Launchkit-org/LaunchKit/gateway/jwt"
	authv1 "github.com/Launchkit-org/LaunchKit/shared/proto/auth/v1"
)

type AuthService interface {
	GetNonce(ctx context.Context, walletAddress string) (nonce string, message string, err error)
	Verify(ctx context.Context, message, signature string, userType string) (*models.UserResponse, *jwt.TokenPayload, error)
	Logout(ctx context.Context, userID, refreshToken string) error
}


type CoreClient interface {
	GetUser(ctx context.Context, id string) (*authv1.User, error)
	GetUserByWallet(ctx context.Context, walletAddress string) (*authv1.User, error)
	CreateUser(ctx context.Context, walletAddress string, userType string) (*authv1.User, error)
	UpdateUser(ctx context.Context, req *authv1.UpdateUserRequest) (*authv1.User, error)
}

func IsUserOnboarded(user *authv1.User) bool {
	if user.UserType == "b2c" {
		return true
	}
	return user.ProjectId != ""
}
