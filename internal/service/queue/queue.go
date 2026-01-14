package queue

import (
	"context"
	"simple_brocker/internal/config"
	"simple_brocker/internal/service/event"
	"simple_brocker/internal/service/queue/dispatch"
	"simple_brocker/internal/service/queue/ingestion"
)

type (
	queue struct {
		ingestion ingestion.Ingestion
		dispatch  dispatch.Dispatch
	}

	Queue interface {
		Producer(ctx context.Context, ch <-chan event.Event)
		Consumer(ctx context.Context, ch chan<- []event.Event)
		Logging(ctx context.Context)
		Close()
	}
)

func New(cfg config.ServiceConf) Queue {
	return &queue{
		ingestion: ingestion.New(cfg),
		dispatch:  dispatch.New(cfg),
	}
}

func (q *queue) Close() {
	q.ingestion.Close()
	q.dispatch.Close()
}

func (q *queue) Producer(ctx context.Context, ch <-chan event.Event) {
	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-ch:
			if !ok {
				return
			}
			q.ingestion.AddEvent(ev)
		}
	}
}

func (q *queue) Consumer(ctx context.Context, ch chan<- []event.Event) {
	for {
		select {
		case <-ctx.Done():
			close(ch)
			return
		default:
			ev := q.dispatch.TakeEvent()
			select {
			case <-ctx.Done():
				close(ch)
				return
			case ch <- ev:
			}
		}
	}
}

func (q *queue) Logging(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ev := q.ingestion.TakeEvent()
			// log.Println("ADD: ", ev.GetData())
			q.dispatch.AddEvent(ev)
		}
	}
}
