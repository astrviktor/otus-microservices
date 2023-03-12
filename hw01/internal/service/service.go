package service

import (
	"context"
	"errors"
	"log"
	"net/http"
	"otus-microservices/hw01/internal/handler"
	"otus-microservices/hw01/internal/middleware"
	"sync"
	"time"
)

type Server struct {
	addr string
	wg   *sync.WaitGroup
	srv  *http.Server
}

func New(addr string) *Server {
	return &Server{
		addr: addr,
		wg:   &sync.WaitGroup{},
		srv:  &http.Server{},
	}
}

func (s *Server) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health/", middleware.Logging((handler.HandleHealth)))

	s.srv = &http.Server{
		Addr:    s.addr,
		Handler: mux,
	}

	log.Println("http server starting on address: " + s.addr)

	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe(): %v", err)
		}
		log.Println("http server stopping")
	}()
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatalf("http server shutdown error: %v", err)
	}

	defer cancel()

	// Wait for ListenAndServe goroutine to close.
	s.wg.Wait()
	log.Println("http server gracefully shutdown")
}
