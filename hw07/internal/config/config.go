package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server  ServerConfig
	Storage StorageConfig
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
