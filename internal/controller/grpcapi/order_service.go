package grpcapi

import (
	"context"
	"gorder-gw/internal/usecase"
	"log"

	orderservice "gorder-gw/internal/generated"
)

type OrderService struct {
	orderservice.UnimplementedOrderServiceServer
	uc *usecase.ConfirmOrder
}

func NewOrderService(uc *usecase.ConfirmOrder) *OrderService {
	return &OrderService{uc: uc}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *orderservice.CreateOrderRequest) (*orderservice.CreateOrderResponse, error) {
	log.Printf("[order-gw] CreateOrder order_id=%s user_id=%s cents=%d curr=%s",
		req.GetOrderId(), req.GetUserId(), req.GetAmountCents(), req.GetCurrency())

	in := usecase.ConfirmOrderInput{
		OrderID:  req.GetOrderId(),
		UserID:   req.GetUserId(),
		Cents:    req.GetAmountCents(),
		Currency: req.GetCurrency(),
	}
	if err := s.uc.Execute(ctx, in); err != nil {
		return nil, err
	}

	return &orderservice.CreateOrderResponse{Status: "ACCEPTED"}, nil
}
