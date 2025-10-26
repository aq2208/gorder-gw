package usecase

import "context"

// EventPublisher is the outbound port to publish events to your bus (Kafka).
type EventPublisher interface {
	PublishOrderSucceeded(ctx context.Context, evt OrderSucceeded) error
}

// OrderSucceeded is the event contract sent to order-api.
type OrderSucceeded struct {
	OrderID  string `json:"orderId"`
	UserID   string `json:"userId"`
	Cents    int64  `json:"cents"`
	Currency string `json:"currency"`
	Status   string `json:"status"` // "SUCCESS"
}
