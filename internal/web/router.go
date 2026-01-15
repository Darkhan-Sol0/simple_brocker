package web

import (
	"simple_brocker/internal/service/most"

	"github.com/labstack/echo/v4"
)

type (
	router struct {
		ioChan most.Most
	}

	Router interface {
		RegisterRoutes(e *echo.Echo)
	}
)

func New(ioChan most.Most) Router {
	return &router{
		ioChan: ioChan,
	}
}

func (r *router) RegisterRoutes(e *echo.Echo) {
	e.POST("/:group", r.RequestEvent)
}
