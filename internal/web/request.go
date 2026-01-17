package web

import (
	"io"
	"net/http"
	"simple_brocker/internal/service/event"

	"github.com/labstack/echo/v4"
)

func (r *router) RequestEvent(ctx echo.Context) error {
	group := ctx.Param("group")
	message, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "error: poka net")
	}
	ev := event.New(group, message)
	r.ioChan.GetIn() <- ev
	return nil
}
