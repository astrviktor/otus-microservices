package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"net/http"
	"otus-microservices/hw06/internal/config"
	"otus-microservices/hw06/internal/server/handlers"
	"otus-microservices/hw06/internal/server/middleware"
	"otus-microservices/hw06/internal/server/prometheus"
	"sync"
	"time"
)

type Server struct {
	addr   string
	wg     *sync.WaitGroup
	server *fasthttp.Server
	log    *zap.Logger
}

//func New(addr string) (*Server, error) {
//	return &Server{
//		addr: addr,
//		wg:   &sync.WaitGroup{},
//		srv:  &http.Server{},
//	}, nil
//}

func New(cfg config.Config, log *zap.Logger) (*Server, error) {
	handler, err := handlers.New(cfg, log)
	if err != nil {
		return nil, err
	}

	prometheus.NewPrometheus()

	r := router.New()

	//r.GET("/metrics", handler.Metrics)
	//r.GET("/metrics", handler.PrometheusHandler())

	r.GET("/health/", middleware.Logging(log, handler.HandleHealth))
	//r.GET("/testing", middleware.Logging(log, handler.HandleTesting))

	r.POST("/login", middleware.Logging(log, handler.Login))
	r.POST("/logout", middleware.Logging(log, handler.Logout))

	r.POST("/profile", middleware.Logging(log, handler.CreateProfile))

	r.GET("/profile/{id}", middleware.Logging(log, handler.Auth(handler.ReadProfile)))
	r.PUT("/profile/{id}", middleware.Logging(log, handler.Auth(handler.UpdateProfile)))
	r.DELETE("/profile/{id}", middleware.Logging(log, handler.Auth(handler.DeleteProfile)))

	srv := &fasthttp.Server{
		Handler: r.Handler,
		//MaxRequestsPerConn: 10,
		//MaxConnsPerIP:      2,
	}

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	return &Server{
		addr:   addr,
		wg:     &sync.WaitGroup{},
		server: srv,
		log:    log,
	}, nil
}

func (s *Server) Start() {
	s.log.Info("http server starting on address: " + s.addr)

	s.wg.Add(1)

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
	s.log.Info("http server gracefully shutdown")
}
