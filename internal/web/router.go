package web

import (
	"context"
	"simple_brocker/internal/config"
	"simple_brocker/internal/service/event"

	"github.com/labstack/echo/v4"
)

type (
	router struct {
		cfg config.Config
		ch  chan event.Event
	}

	Router interface {
		RegisterRoutes(e *echo.Echo)
		ResponseEvent(ctx context.Context, ch chan []event.Event) error
	}
)

func New(cfg config.Config, ch chan event.Event) Router {
	return &router{
		cfg: cfg,
		ch:  ch,
	}
}

func (r *router) RegisterRoutes(e *echo.Echo) {
	e.POST("/:group", r.RequestEvent)
}
