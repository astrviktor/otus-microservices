package main

import (
	"net"
	"os"
	"os/signal"
	"otus-microservices/hw01/internal/service"
	"syscall"
)

func main() {
	addr := net.JoinHostPort("", "8000")

	service := service.New(addr)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	service.Start()

	<-stop

	service.Stop()
}
