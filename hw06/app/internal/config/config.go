package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Server  ServerConfig
	Storage StorageConfig
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type StorageConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Scheme   string `yaml:"scheme"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func NewConfig(path string) Config {
	var config Config

	file, err := os.ReadFile(path)
	if err != nil {
		log.Println(err.Error())
		return DefaultConfig()
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Println(err.Error())
		return DefaultConfig()
	}

	return config
}

func DefaultConfig() Config {
	log.Println("default config")

	return Config{
		Server: ServerConfig{
			Host: "",
			Port: 8000,
		},

		Storage: StorageConfig{
			//Host:     "127.0.0.1",
			//Port:     5432,
			//Scheme:   "users",
			//User:     "user",
			//Password: "password",
		},
	}
}
