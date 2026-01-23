package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"simple_brocker/internal/config"
	"simple_brocker/internal/service/processor"
	"simple_brocker/internal/service/thread"
	"simple_brocker/internal/web/request"
	"simple_brocker/internal/web/response"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

type (
	server struct {
		httpDriver *echo.Echo
		cfg        config.Config
		thread     thread.Thread
		processor  processor.Processor
		request    request.Request
		response   response.Response
	}

	Server interface {
		Run()
	}
)

func New() Server {
	cfg := config.GetConfig()
	thread := thread.New(cfg)
	processor := processor.New(thread)
	request := request.New(cfg, thread.GetIn())
	response := response.New(cfg, thread.GetOut())
	return &server{
		httpDriver: echo.New(),
		cfg:        cfg,
		thread:     thread,
		processor:  processor,
		request:    request,
		response:   response,
	}
}

func (s *server) start() {
	if s.cfg.GetTLS().Enabled {
		if s.cfg.GetTLS().CertPath == "" || s.cfg.GetTLS().KeyPath == "" {
			log.Fatal("TLS enabled but cert_path or key_path not specified")
		}
		certFile := s.cfg.GetTLS().CertPath
		keyFile := s.cfg.GetTLS().KeyPath

		log.Printf("Starting HTTPS server at %s", s.cfg.GetAddress())
		if err := s.httpDriver.StartTLS(s.cfg.GetAddress(), certFile, keyFile); err != nil {
			log.Fatalf("Failed to start HTTPS server: %v", err)
		}
	} else {
		log.Printf("Starting HTTP server at %s", s.cfg.GetAddress())
		if err := s.httpDriver.Start(s.cfg.GetAddress()); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}
}

func (s *server) shutdown(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := s.httpDriver.Shutdown(shutdownCtx); err != nil {
	}
}

func (s *server) Run() {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	defer s.thread.Close()
	s.request.Req(s.httpDriver)

	go s.processor.Producer(ctx)
	go s.processor.Consumer(ctx)
	go s.response.Sender(ctx)

	go s.start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	select {
	case <-ctx.Done():
	case <-quit:
	}

	s.shutdown(ctx)
}
