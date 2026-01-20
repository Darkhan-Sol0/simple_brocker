package dispatch

import (
	"context"
	"simple_brocker/internal/config"
	"simple_brocker/internal/service/event"
	"time"
)

type (
	dispatchQueue struct {
		cfg       config.GroupConf
		queue     chan event.EventIn
		groupName string
	}

	Dispatch interface {
		Close()
		AddEvent(event event.EventIn)
		TakeEvent(ctx context.Context) event.EventOut
	}
)

func New(cfg config.GroupConf, groupName string) Dispatch {
	return &dispatchQueue{
		cfg:       cfg,
		queue:     make(chan event.EventIn, 100),
		groupName: groupName,
	}
}

func (d *dispatchQueue) Close() {
	close(d.queue)
}

func (d *dispatchQueue) AddEvent(event event.EventIn) {
	d.queue <- event
}

func (d *dispatchQueue) TakeEvent(ctx context.Context) event.EventOut {
	ev := event.NewEvOut(d.groupName)

	timer := time.NewTimer(d.cfg.GetCoolDown())
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			if ev.Len() > 0 {
				return ev
			}
		case e := <-d.queue:
			ev.AddData(e.GetData())
			if ev.Len() >= d.cfg.GetServiceBatchSize() {
				return ev
			}
			timer.Reset(d.cfg.GetCoolDown())
		case <-timer.C:
			if ev.Len() > 0 {
				return ev
			}
			timer.Reset(d.cfg.GetCoolDown())
		}
	}
}

func (d *dispatchQueue) GetGroupName() string {
	return d.groupName
}
