package queue

import (
	"context"
	"simple_brocker/internal/config"
	"simple_brocker/internal/service/dispatch"
	"simple_brocker/internal/service/most"
)

type (
	queue struct {
		most     most.Most
		dispatch map[string]dispatch.Dispatch
	}

	Queue interface {
		Producer(ctx context.Context)
		Consumer(ctx context.Context)
		Close()
	}
)

func New(cfg config.Config, most most.Most) Queue {
	ds := make(map[string]dispatch.Dispatch)
	for i := range cfg.GetGroups() {
		ds[i] = dispatch.New(cfg.GetGroup(i))
	}

	return &queue{
		most:     most,
		dispatch: ds,
	}
}

func (q *queue) Close() {
	for _, i := range q.dispatch {
		i.Close()
	}
}

func (q *queue) Producer(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-q.most.GetIn():
			if !ok {
				return
			}
			q.dispatch[ev.GetGroup()].AddEvent(ev)
		}
	}
}

func (q *queue) Consumer(ctx context.Context) {
	for _, i := range q.dispatch {
		go func(d dispatch.Dispatch) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					ev := d.TakeEvent()
					select {
					case <-ctx.Done():
						return
					case q.most.GetOut() <- ev:
					}
				}
			}
		}(i)
	}
}
