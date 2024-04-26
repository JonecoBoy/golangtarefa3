package service

import (
	"context"
	"io"

	"github.com/jonecoboy/golangtarefa3/internal/infra/grpc/pb"
	"github.com/jonecoboy/golangtarefa3/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) CreateOrderStream(stream pb.OrderService_CreateOrderStreamServer) error {
	orders := &pb.CreateOrderResponseList{}
	for {
		order, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(orders)
		}
		dto := usecase.OrderInputDTO{
			ID:    order.Id,
			Price: float64(order.Price),
			Tax:   float64(order.Tax),
		}
		output, err := s.CreateOrderUseCase.Execute(dto)
		if err != nil {
			return err
		}
		orders.CreateOrderResponses = append(orders.CreateOrderResponses, &pb.CreateOrderResponse{
			Id:         output.ID,
			Price:      float32(output.Price),
			Tax:        float32(output.Tax),
			FinalPrice: float32(output.FinalPrice),
		})
	}
}
