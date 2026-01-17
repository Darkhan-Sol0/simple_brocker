package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"simple_brocker/internal/config"
	"simple_brocker/internal/service/most"
	"simple_brocker/internal/service/thread"
	"simple_brocker/internal/web"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

type (
	server struct {
		httpDriver *echo.Echo
		cfg        config.Config
		thread     thread.Thread
		router     web.Router
	}

	Server interface {
		Run()
	}
)

func New() Server {
	cfg := config.GetConfig()
	return &server{
		httpDriver: echo.New(),
		cfg:        cfg,
	}
}

func (s *server) start() {
	log.Printf("Starting server at %s\n", s.cfg.GetAddress())
	if err := s.httpDriver.Start(s.cfg.GetAddress()); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}

func (s *server) shutdown(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	log.Printf("Shutting down the server...")
	if err := s.httpDriver.Shutdown(shutdownCtx); err != nil {
		log.Printf("Graceful shutdown failed with error: %v\n", err)
	}
	log.Println("Server gracefully stopped.")
}

func (s *server) Run() {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	ioChan := most.New(s.cfg)
	defer ioChan.Close()
	s.thread = thread.New(s.cfg, ioChan)
	s.thread.Run(ctx)
	defer s.thread.Close()
	s.router = web.New(s.cfg, ioChan)
	s.router.RegisterRoutes(s.httpDriver)
	go s.start()
	go s.router.ResponseEvent(ctx, ioChan.GetOut())
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("Context canceled, stopping server...")
	case <-quit:
		log.Println("Received termination signal, stopping server...")
	}
	s.shutdown(ctx)
}
