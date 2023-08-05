package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"net/http"
	"otus-microservices/notification/internal/broker"
	"otus-microservices/notification/internal/broker/rabbitmq"
	"otus-microservices/notification/internal/config"
	"otus-microservices/notification/internal/server/handlers"
	"otus-microservices/notification/internal/server/middleware"
	"otus-microservices/notification/internal/server/prometheus"
	"sync"
	"time"
)

type Server struct {
	addr    string
	wg      *sync.WaitGroup
	server  *fasthttp.Server
	handler *handlers.Handler
	broker  broker.InterfaceBroker
	log     *zap.Logger
}

func New(cfg config.Config, log *zap.Logger) (*Server, error) {
	handler, err := handlers.New(cfg, log)
	if err != nil {
		return nil, err
	}

	prometheus.NewPrometheus()

	r := router.New()

	r.GET("/metrics", handler.PrometheusHandler())
	r.GET("/health/", middleware.Logging(log, handler.HandleHealth))

	r.GET("/notification/{id}", middleware.Logging(log, handler.ReadNotification))

	srv := &fasthttp.Server{
		Handler: r.Handler,
	}

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	broker := rabbitmq.New(cfg.Rabbitmq)

	return &Server{
		addr:    addr,
		wg:      &sync.WaitGroup{},
		handler: handler,
		server:  srv,
		log:     log,
		broker:  broker,
	}, nil
}

func (s *Server) Start() {
	s.log.Info("http server starting on address: " + s.addr)

	s.wg.Add(1)

	go func() {
		s.log.Info("broker starting")

		err := s.broker.Connect()
		if err != nil {
			s.log.Fatal("error broker start", zap.Error(err))
		}

		s.handler.ProcessMessages(s.broker)
	}()

	go func() {
		defer s.wg.Done()

		if err := s.server.ListenAndServe(s.addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Fatal("error ListenAndServe()", zap.Error(err))
		}
		s.log.Info("http server stopping")
	}()
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	if err := s.server.ShutdownWithContext(ctx); err != nil {
		s.log.Fatal("http server shutdown error", zap.Error(err))
	}

	defer cancel()

	// Wait for ListenAndServe goroutine to close.
	s.wg.Wait()
	if err := s.broker.Close(); err != nil {
		s.log.Fatal("broker shutdown error", zap.Error(err))
	}

	s.log.Info("http server gracefully shutdown")
}
