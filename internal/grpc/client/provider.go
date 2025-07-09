package grpc

import (
	"context"
	providerpb "top-up-api/proto/provider"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProviderGRPCClient interface {
	Close()
	ProcessOrder(ctx context.Context, req *providerpb.OrderProcessRequest) error
}

type providerGRPCClient struct {
	conn                *grpc.ClientConn
	GrpcProviderService providerpb.ProviderServiceClient
}

var _ ProviderGRPCClient = (*providerGRPCClient)(nil)

func NewProviderGRPCClient(grcpServerUrl string) (*providerGRPCClient, error) {
	// gRPC Client
	conn, err := grpc.NewClient(grcpServerUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {

		return nil, err
	}

	grpcProviderService := providerpb.NewProviderServiceClient(conn)

	return &providerGRPCClient{
		conn:                conn,
		GrpcProviderService: grpcProviderService,
	}, nil
}

func (p *providerGRPCClient) Close() {
	p.conn.Close()
}

func (p *providerGRPCClient) ProcessOrder(ctx context.Context, req *providerpb.OrderProcessRequest) error {
	_, err := p.GrpcProviderService.ProcessOrder(ctx, req)
	return err
}
