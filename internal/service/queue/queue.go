package queue

import (
	"context"
	"simple_brocker/internal/config"
	"simple_brocker/internal/service/event"
	"simple_brocker/internal/service/queue/dispatch"
	"simple_brocker/internal/service/queue/ingestion"
	"sync"
)

type (
	queue struct {
		ingestion ingestion.Ingestion
		dispatch  map[string]dispatch.Dispatch
	}

	Queue interface {
		Producer(ctx context.Context, ch <-chan event.Event)
		Consumer(ctx context.Context, ch chan<- []event.Event)
		Logging(ctx context.Context)
		Close()
	}
)

func New(cfg config.Config) Queue {
	ds := make(map[string]dispatch.Dispatch)
	for i := range cfg.GetGroups() {
		ds[i] = dispatch.New(cfg.GetGroup(i))
	}

	return &queue{
		ingestion: ingestion.New(),
		dispatch:  ds,
	}
}

func (q *queue) Close() {
	q.ingestion.Close()
	for _, i := range q.dispatch {
		i.Close()
	}
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
					ev := d.TakeEvent()
					select {
					case <-ctx.Done():
						return
					case ch <- ev:
					}
				}
			}
		}(i)
	}
	wg.Wait()
	close(ch)
}

func (q *queue) Logging(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ev := q.ingestion.TakeEvent()
			// log.Println("ADD: ", ev.GetData()) // Нужна отдельная функция, для запись. Пока так.
			q.dispatch[ev.GetGroup()].AddEvent(ev)
		}
	}
}
