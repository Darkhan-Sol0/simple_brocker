package web

import (
	"context"
	"simple_brocker/internal/config"
	"simple_brocker/internal/service/event"
	"simple_brocker/internal/service/most"

	"github.com/labstack/echo/v4"
)

type (
	router struct {
		cfg    config.Config
		ioChan most.Most
	}

	Router interface {
		RegisterRoutes(e *echo.Echo)
		ResponseEvent(ctx context.Context, ch <-chan []event.Event) error
	}
)

func New(cfg config.Config, ioChan most.Most) Router {
	return &router{
		cfg:    cfg,
		ioChan: ioChan,
	}
}

func (r *router) RegisterRoutes(e *echo.Echo) {
	e.POST("/:group", r.RequestEvent)
}
