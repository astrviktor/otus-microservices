package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server    ServerConfig
	Storage   StorageConfig
	Payment   PaymentConfig
	Warehouse WarehouseConfig
	Delivery  DeliveryConfig
}

type ServerConfig struct {
	Host string
	Port int `default:"8000"`
}

type StorageConfig struct {
	Host     string `default:"localhost"`
	Port     int    `default:"5432"`
	Scheme   string `default:"orders"`
	User     string `default:"user"`
	Password string `default:"password"`
}

type PaymentConfig struct {
	Host string `default:"localhost"`
	Port int    `default:"8001"`
}

type WarehouseConfig struct {
	Host string `default:"localhost"`
	Port int    `default:"8002"`
}

type DeliveryConfig struct {
	Host string `default:"localhost"`
	Port int    `default:"8003"`
}

func NewConfig() Config {
	var config Config
	return config
}

func PrintUsage(config Config) {
	err := envconfig.Usage("orderservice", &config)
	if err != nil {
		fmt.Printf("Fail to print envconfig usage: %s", err.Error())
	}

	return
}
