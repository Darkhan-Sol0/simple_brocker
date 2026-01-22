package request

import (
	"io"
	"net/http"
	"simple_brocker/internal/service/container"
	"time"

	"github.com/labstack/echo/v4"
)

type (
	request struct {
		chanIn  chan<- container.Container
		sem     chan struct{}
		timeout time.Duration
	}

	Request interface {
		Req(e *echo.Echo)
	}
)

func New(chanIn chan<- container.Container) Request {
	sem := make(chan struct{}, 10)
	for i := 0; i < 10; i++ {
		sem <- struct{}{}
	}

	return &request{
		chanIn:  chanIn,
		sem:     sem,
		timeout: 5 * time.Second,
	}
}

func (r *request) Req(e *echo.Echo) {
	e.POST(":group", r.RequestIn)
}

func (r *request) RequestIn(ctx echo.Context) error {
	select {
	case <-r.sem:
		defer func() { r.sem <- struct{}{} }()

		group := ctx.Param("group")
		data, err := io.ReadAll(ctx.Request().Body)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "failed to read request body: " + err.Error(),
			})
		}
		if len(data) == 0 {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "request body is empty",
			})
		}
		containerData := container.Container{
			Group: group,
			Data:  data,
		}
		select {
		case r.chanIn <- containerData:
			return ctx.JSON(http.StatusOK, map[string]interface{}{
				"message": "request accepted for processing",
				"group":   group,
			})
		case <-time.After(3 * time.Second):
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "request queue is full, try again later",
			})
		case <-ctx.Request().Context().Done():
			return ctx.JSON(http.StatusRequestTimeout, map[string]string{
				"error": "request cancelled while queueing",
			})
		}
	case <-time.After(r.timeout):
		return ctx.JSON(http.StatusServiceUnavailable, map[string]interface{}{
			"error":          "too many concurrent requests",
			"max_concurrent": cap(r.sem),
			"timeout":        r.timeout.String(),
		})
	case <-ctx.Request().Context().Done():
		return ctx.JSON(http.StatusRequestTimeout, map[string]string{
			"error": "request cancelled while waiting for semaphore",
		})
	}
}
