package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"net/http"
	"otus-microservices/orderservice/internal/config"
	"otus-microservices/orderservice/internal/server/handlers"
	"otus-microservices/orderservice/internal/server/middleware"
	"otus-microservices/orderservice/internal/server/prometheus"
	"sync"
	"time"
)

type Server struct {
	addr    string
	wg      *sync.WaitGroup
	server  *fasthttp.Server
	log     *zap.Logger
	handler *handlers.Handler
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

	//r.POST("/order", middleware.Logging(log, handler.CreateOrder))
	r.POST("/order", middleware.Logging(log, handler.CreateOrder))
	r.GET("/order/{id}", middleware.Logging(log, handler.ReadOrder))
	r.PUT("/order/{id}", middleware.Logging(log, handler.UpdateOrder))
	r.DELETE("/order/{id}", middleware.Logging(log, handler.DeleteUser))

	srv := &fasthttp.Server{
		Handler: r.Handler,
	}

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	return &Server{
		addr:    addr,
		wg:      &sync.WaitGroup{},
		server:  srv,
		log:     log,
		handler: handler,
	}, nil
}

func (s *Server) Start() {
	s.log.Info("http server starting on address: " + s.addr)

	s.wg.Add(1)

	go func() {
		s.log.Info("broker starting")

		err := s.handler.Broker.Connect()
		if err != nil {
			s.log.Fatal("error broker start", zap.Error(err))
		}
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

	if err := s.handler.Broker.Close(); err != nil {
		s.log.Fatal("broker shutdown error", zap.Error(err))
	}

	// Wait for ListenAndServe goroutine to close.
	s.wg.Wait()
	s.log.Info("http server gracefully shutdown")
}
