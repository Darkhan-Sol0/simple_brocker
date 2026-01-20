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
		ch  chan event.EventIn
	}

	Router interface {
		RegisterRoutes(e *echo.Echo)
		ResponseEvent(ctx context.Context, ch <-chan event.EventOut)
	}
)

func New(cfg config.Config, ch chan event.EventIn) Router {
	return &router{
		cfg: cfg,
		ch:  ch,
	}
}

func (r *router) RegisterRoutes(e *echo.Echo) {
	e.POST("/:group", r.RequestEvent)
}
