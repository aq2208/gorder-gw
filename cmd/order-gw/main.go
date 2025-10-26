package main

import (
	"gorder-gw/cmd/order-gw/app"
	"gorder-gw/configs"
	"log"
	"os"
)

func main() {
	env := os.Getenv("APP_ENV") // dev | staging | prod
	if env == "" {
		env = "dev"
	}

	cfg, err := configs.Load("configs", env)
	if err != nil {
		log.Fatal(err)
	}

	app, cleanup, err := app.InitWithConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	log.Printf("order-gw app (%s) started", app)
	log.Printf("order-gw (%s) listening on %s", env, cfg.GrpcServer.ListenAddr)
}
