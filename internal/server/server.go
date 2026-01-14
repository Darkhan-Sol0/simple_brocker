package server

import (
	"simple_brocker/internal/config"
	"simple_brocker/internal/service/thread"

	"github.com/labstack/echo/v4"
)

type (
	server struct {
		httpDriver echo.Echo
		cfg        config.Config
		thread     thread.Thread
	}

	Server interface {
	}
)

func New() Server {
	return &server{}
}

func (s *server) start() {

}

func (s *server) shutdown() {

}

func (s *server) Run() {

}
