package event

import (
	"encoding/json"
)

type (
	eventIn struct {
		group string
		data  json.RawMessage
	}
	EventIn interface {
		GetGroup() string
		GetData() json.RawMessage
	}

	eventOut struct {
		group string
		data  []json.RawMessage
	}

	EventOut interface {
		AddData(data json.RawMessage)
		GetGroup() string
		GetData() []json.RawMessage
		Len() int
	}
)

func NewEvIn(group string, data json.RawMessage) EventIn {
	return &eventIn{
		group: group,
		data:  data,
	}
}

func (e *eventIn) GetGroup() string {
	return e.group
}

func (e *eventIn) GetData() json.RawMessage {
	return e.data
}

func NewEvOut(group string) EventOut {
	return &eventOut{
		group: group,
	}
}

func (e *eventOut) AddData(data json.RawMessage) {
	e.data = append(e.data, data)
}

func (e *eventOut) GetGroup() string {
	return e.group
}

func (e *eventOut) GetData() []json.RawMessage {
	return e.data
}

func (e *eventOut) Len() int {
	return len(e.data)
}
