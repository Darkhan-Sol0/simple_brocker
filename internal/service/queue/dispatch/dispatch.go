package dispatch

import (
	"fmt"
	"simple_brocker/internal/config"
	"simple_brocker/internal/service/event"
	"time"
)

type (
	dispatchQueue struct {
		cfg   config.ServiceConf
		queue chan event.Event
	}

	Dispatch interface {
		Close()
		AddEvent(event event.Event)
		TakeEvent() []event.Event
	}
)

func New(cfg config.ServiceConf) Dispatch {
	return &dispatchQueue{
		cfg:   cfg,
		queue: make(chan event.Event, 100),
	}
}

func (d *dispatchQueue) Close() {
	close(d.queue)
}

func (d *dispatchQueue) AddEvent(event event.Event) {
	d.queue <- event
}

func (d *dispatchQueue) TakeEvent() []event.Event {
	ev := make([]event.Event, 0)

	for {
		select {
		case e := <-d.queue:
			ev = append(ev, e)
			if len(ev) >= d.cfg.GetServiceBatchSize() {
				return ev
			}
		case <-time.After(d.cfg.GetCoolDown()):
			if len(ev) > 0 {
				fmt.Println("Zaderzka")
				return ev
			}
			continue
		}
	}
}
