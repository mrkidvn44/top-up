package grpc

import (
	"top-up-api/internal/service"
	pb "top-up-api/proto/order"

	"google.golang.org/grpc"
)

type GRPCServiceContainer struct {
	OrderGRPCServer *OrderGRPCServer
}

func NewGRPCServiceContainer(services *service.Container) *GRPCServiceContainer {
	return &GRPCServiceContainer{
		OrderGRPCServer: NewOrderGRPCServer(services.OrderService),
	}
}

// Register registers all gRPC servers to the given gRPC server.
func (c *GRPCServiceContainer) Register(server grpc.ServiceRegistrar) {
	pb.RegisterOrderServiceServer(server, c.OrderGRPCServer)
}
