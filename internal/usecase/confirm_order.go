package usecase

import (
	"context"
	"errors"

	"gorder-gw/internal/entity"
)

var ErrValidation = errors.New("invalid order payload")

type ConfirmOrderInput struct {
	OrderID  string
	UserID   string
	Cents    int64
	Currency string
}

type ConfirmOrder struct {
	Bus EventPublisher
}

func NewConfirmOrder(bus EventPublisher) *ConfirmOrder {
	return &ConfirmOrder{Bus: bus}
}

// Execute: validate -> create entity -> mark success -> publish event.
func (uc *ConfirmOrder) Execute(ctx context.Context, in ConfirmOrderInput) error {
	if in.OrderID == "" || in.UserID == "" || in.Cents <= 0 || in.Currency == "" {
		return ErrValidation
	}

	ord := domain.Order{
		ID:     in.OrderID,
		UserID: in.UserID,
		Status: domain.StatusProcessing,
		Amount: domain.Money{
			Cents:    in.Cents,
			Currency: in.Currency,
		},
		ItemsJSON: "",
	}
	ord.MarkSuccess() // entity-only status change

	// Publish for order-api to process/update its DB.
	return uc.Bus.PublishOrderSucceeded(ctx, OrderSucceeded{
		OrderID:  ord.ID,
		UserID:   ord.UserID,
		Cents:    ord.Amount.Cents,
		Currency: ord.Amount.Currency,
		Status:   string(ord.Status),
	})
}
