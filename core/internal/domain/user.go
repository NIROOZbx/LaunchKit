package domain

import (
	"context"
	"time"
)

type User struct {
	ID            string    `json:"id"`
	WalletAddress string    `json:"wallet_address"`
	EnsName       string    `json:"ens_name,omitempty"`
	DisplayName   string    `json:"display_name,omitempty"`
	AvatarUrl     string    `json:"avatar_url,omitempty"`
	UserType      string    `json:"user_type"` // 'b2c', 'b2b', 'admin'
	TwitterID     string    `json:"twitter_id,omitempty"`
	TwitterHandle string    `json:"twitter_handle,omitempty"`
	DiscordID     string    `json:"discord_id,omitempty"`
	DiscordHandle string    `json:"discord_handle,omitempty"`
	DiscordToken  string    `json:"discord_token,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	LastSeen      time.Time `json:"last_seen"`
	
	ProjectID     string    `json:"project_id,omitempty"`
	ProjectRole   string    `json:"project_role,omitempty"`
}

type UserService interface {
	CreateUser(ctx context.Context, walletAddress string, userType string) (*User, error)
	GetUser(ctx context.Context, id string) (*User, error)
	GetUserByWallet(ctx context.Context, walletAddress string) (*User, error)
	UpdateUser(ctx context.Context, id string, displayName, avatarUrl, ensName string) (*User, error)
	// LinkTwitter(ctx context.Context, id string, twitterID, twitterHandle string) (*User, error)
	// LinkDiscord(ctx context.Context, id string, discordID, discordHandle, discordToken string) (*User, error)
	// TouchLastSeen(ctx context.Context, id string) error
}

type UserRepository interface {
	Create(ctx context.Context, walletAddress string, ensName, displayName, avatarUrl string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByWallet(ctx context.Context, walletAddress string) (*User, error)
	UpdateProfile(ctx context.Context, id string, displayName, avatarUrl, ensName string) (*User, error)
	// UpsertTwitter(ctx context.Context, id string, twitterID, twitterHandle string) (*User, error)
	// UpsertDiscord(ctx context.Context, id string, discordID, discordHandle, discordToken string) (*User, error)
	// TouchLastSeen(ctx context.Context, id string) error
}
