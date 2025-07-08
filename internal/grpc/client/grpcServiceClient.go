package grpc

import "top-up-api/config"

type GRPCServiceClient struct {
	AuthGRPCClient AuthGRPCClient
}

func NewGRPCServiceClient(
	config config.GrpcClient,
) (*GRPCServiceClient, error) {
	authGRPCClient, err := NewAuthGRPCClient(config.Auth)
	if err != nil {
		return nil, err
	}
	
	return &GRPCServiceClient{
		AuthGRPCClient: authGRPCClient,
	}, nil
}

func (s *GRPCServiceClient) CloseConnection() {
	s.AuthGRPCClient.Close()
}
