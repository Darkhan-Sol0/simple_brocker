package web

import (
	"encoding/json"
	"fmt"
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

	pp := <-r.ioChan.GetOut()
	fmt.Println(pp[0].GetData())
	var text map[string]any
	json.Unmarshal(pp[0].GetData(), &text)
	fmt.Println(text)
	return nil
}
