package grpc

import (
	"top-up-api/internal/service"
	pb "top-up-api/proto/order"

	"google.golang.org/grpc"
)

type GRPCServiceServer struct {
	OrderGRPCServer *OrderGRPCServer
}

func NewGRPCServiceServer(services *service.Container) *GRPCServiceServer {
	return &GRPCServiceServer{
		OrderGRPCServer: NewOrderGRPCServer(services.OrderService),
	}
}

// Register registers all gRPC servers to the given gRPC server.
func (c *GRPCServiceServer) Register(server grpc.ServiceRegistrar) {
	pb.RegisterOrderServiceServer(server, c.OrderGRPCServer)
}
