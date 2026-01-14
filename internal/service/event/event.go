package event

import (
	"encoding/json"
	"time"
)

type (
	event struct {
		data      json.RawMessage
		create_at time.Time
		attempts  int
	}
	Event interface {
		GetData() json.RawMessage
		GetCreateDate() time.Time
		GetCreateDateF() string
		GetAttempts() int
		MakeAttempts()
	}
)

func New(data json.RawMessage) Event {
	return &event{
		data:      data,
		create_at: time.Now(),
		attempts:  0,
	}
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
