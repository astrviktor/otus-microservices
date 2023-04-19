package main

import (
	"flag"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"otus-microservices/hw03/internal/config"
	"otus-microservices/hw03/internal/logger"
	"otus-microservices/hw03/internal/service"
	"syscall"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config.yaml", "Path to configuration file")
}

func main() {
	log, err := logger.New()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	flag.Parse()

	cfg := config.NewConfig(configFile)

	err = envconfig.Process("crudservice", &cfg)
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
