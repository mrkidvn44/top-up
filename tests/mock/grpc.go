package mock

import (
	"context"

	"top-up-api/internal/model"
	pb "top-up-api/proto/provider"

	"github.com/stretchr/testify/mock"
)

// ProviderGRPCClientMock mocks the provider GRPC client
type ProviderGRPCClientMock struct {
	mock.Mock
}

func (m *ProviderGRPCClientMock) ProcessOrder(ctx context.Context, req *pb.OrderProcessRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *ProviderGRPCClientMock) Close() {
	m.Called()
}

// GRPCServiceClientMock mocks the GRPC service client struct
type GRPCServiceClientMock struct {
	mock.Mock
	ProviderGRPCClients map[string]*ProviderGRPCClientMock
}

func (m *GRPCServiceClientMock) BuildProviderGRPCClients(providers []model.Provider) {
	m.Called(providers)
}

func (m *GRPCServiceClientMock) CloseConnection() {
	m.Called()
}
