package app

import (
	"context"
	"gorder-gw/configs"
	"log"
	"os/signal"
	"syscall"
	"time"

	"gorder-gw/internal/controller/grpcapi"
)

type App struct {
}

func InitWithConfig(cfg configs.Config) error {
	// init logger
	//logger, _ := observ.NewLogger()
	//defer logger.Sync()

	// init context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	svc := grpcapi.NewOrderService()

	// trap SIGINT/SIGTERM for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := RunGRPC(ctx, cfg, svc); err != nil {
		log.Fatalf("gRPC server stopped: %v", err)
	}
	log.Println("gRPC server exited")

	return nil
}
