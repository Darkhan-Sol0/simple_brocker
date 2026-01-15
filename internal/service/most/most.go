package most

import (
	"simple_brocker/internal/config"
	"simple_brocker/internal/service/event"
)

type (
	ioChan struct {
		cfg   config.Config
		chIn  map[string]chan event.Event
		chOut map[string]chan []event.Event
	}

	Most interface {
		MakeMost()
		Close()

		GetIn() map[string]chan event.Event
		GetOut() map[string]chan []event.Event
		GetInChan(key string) chan event.Event
		GetOutChan(key string) chan []event.Event
	}
)

func New(cfg config.Config) Most {
	return &ioChan{
		cfg:   cfg,
		chIn:  make(map[string]chan event.Event),
		chOut: make(map[string]chan []event.Event),
	}
}

func (io *ioChan) MakeMost() {
	for i := range io.cfg.GetServices() {
		io.chIn[i] = make(chan event.Event, 100)
		io.chOut[i] = make(chan []event.Event, 10)
	}
}

func (io *ioChan) GetIn() map[string]chan event.Event {
	return io.chIn
}

func (io *ioChan) GetOut() map[string]chan []event.Event {
	return io.chOut
}

func (io *ioChan) GetInChan(key string) chan event.Event {
	return io.chIn[key]
}

func (io *ioChan) GetOutChan(key string) chan []event.Event {
	return io.chOut[key]
}

func (io *ioChan) Close() {
	for _, i := range io.chIn {
		close(i)
	}

	for _, i := range io.chOut {
		close(i)
	}
}
