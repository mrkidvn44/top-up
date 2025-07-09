package grpc

import "top-up-api/config"

type GRPCServiceClient struct {
	AuthGRPCClient     AuthGRPCClient
	ProviderGRPCClient ProviderGRPCClient
}

func NewGRPCServiceClient(
	config config.GrpcClient,
) (*GRPCServiceClient, error) {
	authGRPCClient, err := NewAuthGRPCClient(config.Auth)
	if err != nil {
		return nil, err
	}

	providerGRPCClient, err := NewProviderGRPCClient(config.Provider)
	if err != nil {
		return nil, err
	}
	return &GRPCServiceClient{
		AuthGRPCClient:     authGRPCClient,
		ProviderGRPCClient: providerGRPCClient,
	}, nil
}

func (s *GRPCServiceClient) CloseConnection() {
	s.AuthGRPCClient.Close()
	s.ProviderGRPCClient.Close()
}
