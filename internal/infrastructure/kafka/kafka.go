package kafka

import (
	"time"

	"github.com/IBM/sarama"
)

func MustKafkaSyncProducer(brokers []string, clientID string) sarama.SyncProducer {
	cfg := sarama.NewConfig()

	// Required for idempotent producer
	cfg.Producer.Idempotent = true
	cfg.Net.MaxOpenRequests = 1
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 10
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Idempotent = true
	cfg.Net.DialTimeout = 5 * time.Second
	cfg.ClientID = clientID
	prod, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		panic(err)
	}
	return prod
}
