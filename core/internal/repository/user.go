package repository

import (
	"context"
	"fmt"

	"github.com/Launchkit-org/LaunchKit/core/internal/domain"
	"github.com/Launchkit-org/LaunchKit/core/internal/pgconv"
	db "github.com/Launchkit-org/LaunchKit/db/sqlc"
)

type userRepository struct {
	queries *db.Queries
}

func NewUserRepository(queries *db.Queries) domain.UserRepository {
	return &userRepository{
		queries: queries,
	}
}

func (r *userRepository) Create(ctx context.Context, walletAddress string, ensName, displayName, avatarUrl string) (*domain.User, error) {
	dbUser, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		WalletAddress: walletAddress,
		EnsName:       pgconv.StringToText(ensName),
		DisplayName:   pgconv.StringToText(displayName),
		AvatarUrl:     pgconv.StringToText(avatarUrl),
	})
	if err != nil {
		return nil, fmt.Errorf("userRepository.Create: %w", err)
	}
	return mapToDomainUser(dbUser), nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	uuidVal := pgconv.StringToUUID(id)
	if !uuidVal.Valid {
		return nil, fmt.Errorf("userRepository.GetByID: invalid UUID: %s", id)
	}
	dbUser, err := r.queries.GetUserByID(ctx, uuidVal)
	if err != nil {
		return nil, fmt.Errorf("userRepository.GetByID: %w", err)
	}
	return mapToDomainUser(dbUser), nil
}

func (r *userRepository) GetByWallet(ctx context.Context, walletAddress string) (*domain.User, error) {
	dbUser, err := r.queries.GetUserByWallet(ctx, walletAddress)
	if err != nil {
		return nil, fmt.Errorf("userRepository.GetByWallet: %w", err)
	}
	return mapToDomainUser(dbUser), nil
}

func (r *userRepository) UpdateProfile(ctx context.Context, id string, displayName, avatarUrl, ensName string) (*domain.User, error) {
	uuidVal := pgconv.StringToUUID(id)
	if !uuidVal.Valid {
		return nil, fmt.Errorf("userRepository.UpdateProfile: invalid UUID: %s", id)
	}
	dbUser, err := r.queries.UpdateUserProfile(ctx, db.UpdateUserProfileParams{
		ID:          uuidVal,
		DisplayName: pgconv.StringToText(displayName),
		AvatarUrl:   pgconv.StringToText(avatarUrl),
		EnsName:     pgconv.StringToText(ensName),
	})
	if err != nil {
		return nil, fmt.Errorf("userRepository.UpdateProfile: %w", err)
	}
	return mapToDomainUser(dbUser), nil
}

func mapToDomainUser(u db.User) *domain.User {
	return &domain.User{
		ID:            pgconv.UUIDToString(u.ID),
		WalletAddress: u.WalletAddress,
		EnsName:       pgconv.TextToString(u.EnsName),
		DisplayName:   pgconv.TextToString(u.DisplayName),
		AvatarUrl:     pgconv.TextToString(u.AvatarUrl),
		UserType:      u.UserType,
		TwitterID:     pgconv.TextToString(u.TwitterID),
		TwitterHandle: pgconv.TextToString(u.TwitterHandle),
		DiscordID:     pgconv.TextToString(u.DiscordID),
		DiscordHandle: pgconv.TextToString(u.DiscordHandle),
		DiscordToken:  pgconv.TextToString(u.DiscordToken),
		CreatedAt:     u.CreatedAt.Time,
		UpdatedAt:     u.UpdatedAt.Time,
		LastSeen:      u.LastSeen.Time,
	}
}
