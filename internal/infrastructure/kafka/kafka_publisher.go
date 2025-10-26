package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"gorder-gw/internal/usecase"

	"github.com/IBM/sarama"
)

const DefaultTopic = "order.status.changed"

type KafkaPublisher struct {
	prod  sarama.SyncProducer
	topic string
}

func NewKafkaPublisher(prod sarama.SyncProducer, topic string) *KafkaPublisher {
	if topic == "" {
		topic = DefaultTopic
	}
	return &KafkaPublisher{prod: prod, topic: topic}
}

func (p *KafkaPublisher) PublishOrderSucceeded(ctx context.Context, evt usecase.OrderSucceeded) error {
	b, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		// Use orderID as key to keep partitioning stable
		Key:   sarama.StringEncoder(evt.OrderID),
		Value: sarama.ByteEncoder(b),
	}
	_, _, err = p.prod.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("kafka send: %w", err)
	}
	return nil
}
