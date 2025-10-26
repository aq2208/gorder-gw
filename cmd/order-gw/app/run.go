package app

import (
	"context"
	"gorder-gw/configs"
	"gorder-gw/internal/infrastructure/kafka"
	"gorder-gw/internal/usecase"
	"log"
	"os/signal"
	"syscall"
	"time"

	"gorder-gw/internal/controller/grpcapi"

	"github.com/IBM/sarama"
)

type App struct {
	Kafka sarama.SyncProducer
}

func InitWithConfig(cfg configs.Config) (*App, func(), error) {
	// init logger
	//logger, _ := observ.NewLogger()
	//defer logger.Sync()

	// init context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// init kafka
	prod := kafka.MustKafkaSyncProducer(cfg.KafkaBroker.KafkaBrokers, "order-gw")
	bus := kafka.NewKafkaPublisher(prod, cfg.KafkaBroker.KafkaTopic)
	uc := usecase.NewConfirmOrder(bus)
	svc := grpcapi.NewOrderService(uc)

	// trap SIGINT/SIGTERM for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := RunGRPC(ctx, cfg, svc); err != nil {
		log.Fatalf("gRPC server stopped: %v", err)
	}
	log.Println("gRPC server exited")

	cleanup := func() {
		err := prod.Close()
		if err != nil {
			return
		}
	}

	return &App{Kafka: prod}, cleanup, nil
}
