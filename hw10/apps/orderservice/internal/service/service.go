package service

import (
	"go.uber.org/zap"
	"otus-microservices/orderservice/internal/config"
	"otus-microservices/orderservice/internal/server"
)

type Service struct {
	Config config.Config
	Server *server.Server
}

func New(cfg config.Config, log *zap.Logger) (*Service, error) {
	srv, err := server.New(cfg, log)

	if err != nil {
		return nil, err
	}

	return &Service{
		Config: cfg,
		Server: srv,
	}, nil
}

func (s *Service) Start() {
	s.Server.Start()
}

func (s *Service) Stop() {
	s.Server.Stop()
}
