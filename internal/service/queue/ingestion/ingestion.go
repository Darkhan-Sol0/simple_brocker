package ingestion

import (
	"simple_brocker/internal/service/event"
)

type (
	ingestionQueue struct {
		queue chan event.Event
	}

	Ingestion interface {
		Close()
		AddEvent(event event.Event)
		TakeEvent() event.Event
	}
)

func New() Ingestion {
	return &ingestionQueue{
		queue: make(chan event.Event, 100),
	}
}

func (i *ingestionQueue) Close() {
	close(i.queue)
}

func (i *ingestionQueue) AddEvent(event event.Event) {
	i.queue <- event
}

func (i *ingestionQueue) TakeEvent() event.Event {
	return <-i.queue
}
