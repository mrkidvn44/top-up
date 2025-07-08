package grpc

import (
	"context"
	authpb "top-up-api/proto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IAuthGRPCClient interface {
	Close()
	AuthenticateService(ctx context.Context, req *authpb.AuthenticateServiceRequest) error
}

type AuthGRPCClient struct {
	conn            *grpc.ClientConn
	GrpcAuthService authpb.AuthServiceClient
}

var _ IAuthGRPCClient = (*AuthGRPCClient)(nil)

func NewAuthGRPCClient(grcpServerUrl string) (*AuthGRPCClient, error) {
	// gRPC Client
	conn, err := grpc.NewClient(grcpServerUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {

		return nil, err
	}

	grpcAuthService := authpb.NewAuthServiceClient(conn)
	
	return &AuthGRPCClient{
		conn:            conn,
		GrpcAuthService: grpcAuthService,
	}, nil
}

func (a *AuthGRPCClient) Close() {
	a.conn.Close()
}

func (a *AuthGRPCClient) AuthenticateService(ctx context.Context, req *authpb.AuthenticateServiceRequest) error {
	_, err := a.GrpcAuthService.AuthenticateService(ctx, req)
	return err
}
