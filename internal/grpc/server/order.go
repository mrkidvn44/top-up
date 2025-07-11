package grpc

import (
	"context"
	"fmt"
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
		return &pb.ConfirmOrderResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}
	return &pb.ConfirmOrderResponse{
		Success: true,
		Error:   "",
	}, nil

}

func (s *OrderGRPCServer) UpdateOrderStatus(ctx context.Context, req *pb.OrderUpdateRequest) (*pb.OrderUpdateResponse, error) {
	orderUpdateRequest := mapper.OrderUpdateRequestFromProto(req)
	err := s.orderService.UpdateOrderStatus(ctx, *orderUpdateRequest)
	if err != nil {
		fmt.Print(err)
		return &pb.OrderUpdateResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}
	return &pb.OrderUpdateResponse{
		Success: true,
		Error:   "",
	}, nil
}
