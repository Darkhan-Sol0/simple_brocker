package most

import (
	"simple_brocker/internal/config"
	"simple_brocker/internal/service/event"
)

type (
	ioChan struct {
		cfg   config.Config
		chIn  chan event.Event
		chOut chan []event.Event
	}

	Most interface {
		Close()

		GetIn() chan event.Event
		GetOut() chan []event.Event
	}
)

func New(cfg config.Config) Most {
	return &ioChan{
		cfg:   cfg,
		chIn:  make(chan event.Event, 100),
		chOut: make(chan []event.Event, 10),
	}
}

func (io *ioChan) GetIn() chan event.Event {
	return io.chIn
}

func (io *ioChan) GetOut() chan []event.Event {
	return io.chOut
}

func (io *ioChan) Close() {
	close(io.chIn)
	close(io.chOut)
}
