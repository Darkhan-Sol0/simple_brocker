package dispatch

import (
	"context"
	"fmt"
	"simple_brocker/internal/config"
	"simple_brocker/internal/service/event"
	"time"
)

type (
	dispatchQueue struct {
		cfg       config.GroupConf
		queue     chan event.Event
		groupName string
	}

	Dispatch interface {
		Close()
		AddEvent(event event.Event)
		TakeEvent(ctx context.Context) []event.Event
	}
)

func New(cfg config.GroupConf, groupName string) Dispatch {
	return &dispatchQueue{
		cfg:       cfg,
		queue:     make(chan event.Event, 100),
		groupName: groupName,
	}
}

func (d *dispatchQueue) Close() {
	close(d.queue)
}

func (d *dispatchQueue) AddEvent(event event.Event) {
	d.queue <- event
}

func (d *dispatchQueue) TakeEvent(ctx context.Context) []event.Event {
	ev := make([]event.Event, 0)

	timer := time.NewTimer(d.cfg.GetCoolDown())
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			if len(ev) > 0 {
				return ev
			}
			return nil
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

func (d *dispatchQueue) GetGroupName() string {
	return d.groupName
}
