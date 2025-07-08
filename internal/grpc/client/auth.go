package grpc

import (
	"context"
	authpb "top-up-api/proto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthGRPCClient interface {
	Close()
	AuthenticateService(ctx context.Context, req *authpb.AuthenticateServiceRequest) error
}

type authGRPCClient struct {
	conn            *grpc.ClientConn
	GrpcAuthService authpb.AuthServiceClient
}

var _ AuthGRPCClient = (*authGRPCClient)(nil)

func NewAuthGRPCClient(grcpServerUrl string) (*authGRPCClient, error) {
	// gRPC Client
	conn, err := grpc.NewClient(grcpServerUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {

		return nil, err
	}

	grpcAuthService := authpb.NewAuthServiceClient(conn)

	return &authGRPCClient{
		conn:            conn,
		GrpcAuthService: grpcAuthService,
	}, nil
}

func (a *authGRPCClient) Close() {
	a.conn.Close()
}

func (a *authGRPCClient) AuthenticateService(ctx context.Context, req *authpb.AuthenticateServiceRequest) error {
	_, err := a.GrpcAuthService.AuthenticateService(ctx, req)
	return err
}
