package event

import (
	"encoding/json"
	"time"
)

type (
	event struct {
		group     string
		data      json.RawMessage
		create_at time.Time
		attempts  int
	}
	Event interface {
		GetGroup() string
		GetData() json.RawMessage
		GetCreateDate() time.Time
		GetCreateDateF() string
		GetAttempts() int
		MakeAttempts()
	}
)

func New(group string, data json.RawMessage) Event {
	return &event{
		group:     group,
		data:      data,
		create_at: time.Now(),
		attempts:  0,
	}
}

func (e *event) GetGroup() string {
	return e.group
}

func (e *event) GetData() json.RawMessage {
	return e.data
}

func (e *event) GetCreateDate() time.Time {
	return e.create_at
}

func (e *event) GetCreateDateF() string {
	return e.create_at.Format("02-01-2006")
}

func (e *event) GetAttempts() int {
	return e.attempts
}

func (e *event) MakeAttempts() {
	e.attempts++
}
