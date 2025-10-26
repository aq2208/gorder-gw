package grpcapi

import (
	"context"
	"log"

	orderservice "gorder-gw/internal/generated"
)

type OrderService struct {
	orderservice.UnimplementedOrderServiceServer
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *orderservice.CreateOrderRequest) (*orderservice.CreateOrderResponse, error) {
	// TODO: real work (validate, DB, etc). For now, accept it.
	log.Printf("[order-gw] CreateOrder order_id=%s user_id=%s cents=%d curr=%s",
		req.GetOrderId(), req.GetUserId(), req.GetAmountCents(), req.GetCurrency())

	// reply something meaningful
	return &orderservice.CreateOrderResponse{Status: "ACCEPTED"}, nil
}
