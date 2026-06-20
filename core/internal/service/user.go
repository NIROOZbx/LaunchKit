package service

import (
	"context"

	"github.com/Launchkit-org/LaunchKit/core/internal/domain"
)

type userService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) domain.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, walletAddress string, userType string) (*domain.User, error) {
	// 1. Check if user already exists
	user, err := s.userRepo.GetByWallet(ctx, walletAddress)
	if err == nil && user != nil {

		return user, nil
	}

	// 2. If not, create a new user
	// Note: We use the repository interface which calls the sqlc Create query.
	newUser, err := s.userRepo.Create(ctx, walletAddress, "", "", "")
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) GetUserByWallet(ctx context.Context, walletAddress string) (*domain.User, error) {
	return s.userRepo.GetByWallet(ctx, walletAddress)
}

func (s *userService) UpdateUser(ctx context.Context, id string, displayName, avatarUrl, ensName string) (*domain.User, error) {
	return s.userRepo.UpdateProfile(ctx, id, displayName, avatarUrl, ensName)
}
