package thread

import (
	"simple_brocker/internal/config"
	"simple_brocker/internal/service/container"
)

type (
	thread struct {
		chanIn  chan container.Container
		chanOut map[string]chan container.Container
	}

	Thread interface {
		Close()

		GetIn() chan container.Container
		GetOut() map[string]chan container.Container
	}
)

func New(cfg config.Config) Thread {

	chanIn := make(chan container.Container, 100)
	chanOut := make(map[string]chan container.Container)
	for i := range cfg.GetGroups() {
		chanOut[i] = make(chan container.Container, 100)
	}

	return &thread{
		chanIn:  chanIn,
		chanOut: chanOut,
	}
}

func (t *thread) Close() {
	close(t.chanIn)
	for i := range t.chanOut {
		close(t.chanOut[i])
	}

}

func (t *thread) GetIn() chan container.Container {
	return t.chanIn
}

func (t *thread) GetOut() map[string]chan container.Container {
	return t.chanOut
}
