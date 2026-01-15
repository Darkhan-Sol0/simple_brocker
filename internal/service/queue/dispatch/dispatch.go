package dispatch

import (
	"fmt"
	"simple_brocker/internal/config"
	"simple_brocker/internal/service/event"
	"time"
)

type (
	dispatchQueue struct {
		cfg   config.GroupConf
		queue chan event.Event
	}

	Dispatch interface {
		Close()
		AddEvent(event event.Event)
		TakeEvent() []event.Event
	}
)

func New(cfg config.GroupConf) Dispatch {
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
	timer := time.NewTimer(d.cfg.GetCoolDown())
	defer timer.Stop()

	for {
		select {
		case e := <-d.queue:
			ev = append(ev, e)
			if len(ev) >= d.cfg.GetServiceBatchSize() {
				return ev
			}
			timer.Reset(d.cfg.GetCoolDown())
		case <-timer.C:
			if len(ev) > 0 {
				fmt.Println("Zaderzka")
				return ev
			}
			timer.Reset(d.cfg.GetCoolDown())
		}
	}
}
