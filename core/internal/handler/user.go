package handler

import (
	"context"

	"github.com/Launchkit-org/LaunchKit/core/internal/domain"
	authv1 "github.com/Launchkit-org/LaunchKit/shared/proto/auth/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserHandler struct {
	authv1.UnimplementedUserServiceServer
	userService domain.UserService
}

func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *authv1.CreateUserRequest) (*authv1.UserResponse, error) {
	user, err := h.userService.CreateUser(ctx, req.WalletAddress, req.UserType)
	if err != nil {
		return nil, err
	}

	return &authv1.UserResponse{
		User: toProtoUser(user),
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *authv1.GetUserRequest) (*authv1.UserResponse, error) {
	user, err := h.userService.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &authv1.UserResponse{
		User: toProtoUser(user),
	}, nil
}

func (h *UserHandler) GetUserWallet(ctx context.Context, req *authv1.GetUserByWalletRequest) (*authv1.UserResponse, error) {
	user, err := h.userService.GetUserByWallet(ctx, req.WalletAddress)
	if err != nil {
		return nil, err
	}

	return &authv1.UserResponse{
		User: toProtoUser(user),
	}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *authv1.UpdateUserRequest) (*authv1.UserResponse, error) {
	user, err := h.userService.UpdateUser(ctx, req.Id, req.DisplayName, req.AvatarUrl, req.EnsName)
	if err != nil {
		return nil, err
	}

	return &authv1.UserResponse{
		User: toProtoUser(user),
	}, nil
}

func toProtoUser(u *domain.User) *authv1.User {
	if u == nil {
		return nil
	}

	return &authv1.User{
		Id:            u.ID,
		WalletAddress: u.WalletAddress,
		EnsName:       u.EnsName,
		DisplayName:   u.DisplayName,
		AvatarUrl:     u.AvatarUrl,
		UserType:      u.UserType,
		TwitterId:     u.TwitterID,
		TwitterHandle: u.TwitterHandle,
		DiscordId:     u.DiscordID,
		DiscordHandle: u.DiscordHandle,
		DiscordToken:  u.DiscordToken,
		ProjectId:     u.ProjectID,
		ProjectRole:   u.ProjectRole,
		CreatedAt:     timestamppb.New(u.CreatedAt),
		UpdatedAt:     timestamppb.New(u.UpdatedAt),
		LastSeen:      timestamppb.New(u.LastSeen),
	}
}