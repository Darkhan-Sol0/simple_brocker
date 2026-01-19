package queue

import (
	"context"
	"log"
	"simple_brocker/internal/config"
	"simple_brocker/internal/service/dispatch"
	"simple_brocker/internal/service/event"
	"simple_brocker/internal/service/logging"
	"sync"
)

type (
	queue struct {
		chIn     chan event.Event
		chOut    chan []event.Event
		dispatch map[string]dispatch.Dispatch

		log logging.Logger
	}

	Queue interface {
		Producer(ctx context.Context)
		Consumer(ctx context.Context)
		Close()

		GetIn() chan event.Event
		GetOut() chan []event.Event
	}
)

func (q *queue) GetIn() chan event.Event {
	return q.chIn
}

func (q *queue) GetOut() chan []event.Event {
	return q.chOut
}

func New(cfg config.Config) Queue {
	ds := make(map[string]dispatch.Dispatch)
	for i := range cfg.GetGroups() {
		ds[i] = dispatch.New(cfg.GetGroup(i), i)
	}

	l, _ := logging.New()

	return &queue{
		chIn:     make(chan event.Event, 100),
		chOut:    make(chan []event.Event, 10),
		dispatch: ds,

		log: l,
	}
}

func (q *queue) Close() {
	for _, i := range q.dispatch {
		i.Close()
	}
	close(q.chIn)
	close(q.chOut)
	q.log.Close()
}

func (q *queue) Producer(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-q.chIn:
			if !ok {
				return
			}
			d, exist := q.dispatch[ev.GetGroup()]
			if !exist {
				log.Printf("group %s not found", ev.GetGroup())
				continue
			}
			q.log.LogEvent(ev)
			d.AddEvent(ev)
		}
	}
}

func (q *queue) Consumer(ctx context.Context) {
	var wg sync.WaitGroup
	for _, i := range q.dispatch {
		wg.Add(1)
		go func(d dispatch.Dispatch) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					ev := d.TakeEvent(ctx)
					select {
					case <-ctx.Done():
						return
					case q.chOut <- ev:
					}
				}
			}
		}(i)
	}
	wg.Wait()
}
