package client

import (
	"context"
	"fmt"

	authv1 "github.com/Launchkit-org/LaunchKit/shared/proto/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CoreClient struct {
	conn *grpc.ClientConn
	grpcClient authv1.UserServiceClient
}

func NewCoreClient(addr string)(*CoreClient,error){

	conn,err:=grpc.NewClient(addr,grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial core gRPC server: %w", err)
	}

	return &CoreClient{
		conn:   conn,
		grpcClient: authv1.NewUserServiceClient(conn),
	}, nil
}

func (c *CoreClient) Close() error {
	return c.conn.Close()
}


func (c *CoreClient) GetUser(ctx context.Context, id string) (*authv1.User, error) {
	resp, err := c.grpcClient.GetUser(ctx, &authv1.GetUserRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return resp.User, nil
}

func (c *CoreClient) GetUserByWallet(ctx context.Context, walletAddress string) (*authv1.User, error) {
	resp, err := c.grpcClient.GetUserWallet(ctx, &authv1.GetUserByWalletRequest{
		WalletAddress: walletAddress,
	})
	if err != nil {
		return nil, err 
	}
	return resp.User, nil
}

func (c *CoreClient) CreateUser(ctx context.Context, walletAddress string,userType string) (*authv1.User, error) {
	resp, err := c.grpcClient.CreateUser(ctx, &authv1.CreateUserRequest{
		WalletAddress: walletAddress,
		UserType:     userType, 
	})
	if err != nil {
		return nil, err
	}
	return resp.User, nil
}

func (c *CoreClient) UpdateUser(ctx context.Context, req *authv1.UpdateUserRequest) (*authv1.User, error) {
	resp, err := c.grpcClient.UpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.User, nil
}
