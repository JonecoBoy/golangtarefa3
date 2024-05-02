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
	GetOrderUseCase    usecase.GetOrderUseCase
	ListOrdersUseCase  usecase.ListOrdersUseCase
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase, getOrderUseCase usecase.GetOrderUseCase, listOrdersUseCase usecase.ListOrdersUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		GetOrderUseCase:    getOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.OrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) CreateOrderStream(stream pb.OrderService_CreateOrderStreamServer) error {
	orders := &pb.OrderResponseList{}
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
		orders.CreateOrderResponses = append(orders.CreateOrderResponses, &pb.OrderResponse{
			Id:         output.ID,
			Price:      float32(output.Price),
			Tax:        float32(output.Tax),
			FinalPrice: float32(output.FinalPrice),
		})
	}
}

func (s *OrderService) GetOrder(ctx context.Context, in *pb.GetOrderRequest) (*pb.OrderResponse, error) {
	// Convert the GetOrderRequest to the input DTO for your use case

	// Execute the use case with the input DTO
	output, err := s.GetOrderUseCase.Execute(in.Id)
	if err != nil {
		return nil, err
	}

	// Convert the output DTO to the OrderResponse
	return &pb.OrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, in *pb.ListOrderRequest) (*pb.OrderResponseList, error) {
	// Execute the use case
	outputs, err := s.ListOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}

	// Convert the output DTOs to the OrderResponseList
	responses := make([]*pb.OrderResponse, len(outputs))
	for i, output := range outputs {
		responses[i] = &pb.OrderResponse{
			Id:         output.ID,
			Price:      float32(output.Price),
			Tax:        float32(output.Tax),
			FinalPrice: float32(output.FinalPrice),
		}
	}

	return &pb.OrderResponseList{
		CreateOrderResponses: responses,
	}, nil
}
