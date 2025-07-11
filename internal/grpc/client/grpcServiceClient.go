package grpc

import (
	"top-up-api/config"
	"top-up-api/internal/model"
)

type GRPCServiceClient struct {
	AuthGRPCClient      AuthGRPCClient
	ProviderGRPCClients map[string]ProviderGRPCClient
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

func (s *GRPCServiceClient) BuildProviderGRPCClients(providers []model.Provider) {
	clients := make(map[string]ProviderGRPCClient)
	for _, provider := range providers {
		if provider.Type == "grpc" {
			client, err := NewProviderGRPCClient(provider.Source)
			if err != nil {
				panic("failed to create gRPC client for " + provider.Code + ": " + err.Error())
			}
			clients[provider.Code] = client
		}
	}
	s.ProviderGRPCClients = clients
}

func (s *GRPCServiceClient) CloseConnection() {
	s.AuthGRPCClient.Close()
	for _, service := range s.ProviderGRPCClients {
		service.Close()
	}

}
