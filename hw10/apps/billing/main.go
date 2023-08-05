package main

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"otus-microservices/billing/internal/config"
	"otus-microservices/billing/internal/logger"
	"otus-microservices/billing/internal/service"
	"syscall"
)

func main() {
	log, err := logger.New()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	cfg := config.NewConfig()
	if len(os.Args) > 1 {
		config.PrintUsage(cfg)
		return
	}

	err = envconfig.Process("billing", &cfg)
	if err != nil {
		log.Fatal("fail to get env", zap.Error(err))
	}

	service, err := service.New(cfg, log)
	if err != nil {
		log.Fatal("fail to create service", zap.Error(err))
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	service.Start()

	<-stop

	service.Stop()
}
