package grpc

import (
	"context"
	"top-up-api/internal/mapper"
	"top-up-api/internal/service"
	pb "top-up-api/proto/order"
)

type OrderGRPCServer struct {
	pb.UnimplementedOrderServiceServer
	orderService service.OrderService
}

func NewOrderGRPCServer(orderService service.OrderService) *OrderGRPCServer {
	return &OrderGRPCServer{
		orderService: orderService,
	}
}

func (s *OrderGRPCServer) ConfirmOrder(ctx context.Context, req *pb.OrderConfirmRequest) (*pb.ConfirmOrderResponse, error) {
	orderConfirmRequest := mapper.OrderConfirmRequestFromProto(req)
	err := s.orderService.ConfirmOrder(ctx, *orderConfirmRequest)
	if err != nil {
		return nil, err
	}
	return &pb.ConfirmOrderResponse{
		Success: true,
		Error:   "",
	}, nil

}
