package processor

import (
	"context"
	"simple_brocker/internal/config"
	"simple_brocker/internal/service/container"
	"simple_brocker/internal/service/fsaver"
	"simple_brocker/internal/service/thread"
)

type (
	processor struct {
		thread thread.Thread
		fsaver fsaver.Fsaver
	}

	Processor interface {
		Producer(ctx context.Context)
		Consumer(ctx context.Context)
	}
)

func New(thread thread.Thread) Processor {
	return &processor{
		thread: thread,
		fsaver: fsaver.New(config.GetConfig()),
	}
}

func (p *processor) Producer(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-p.thread.GetIn():
			p.fsaver.LogData(data)
		}
	}

}

func (p *processor) Consumer(ctx context.Context) {
	for i, j := range p.thread.GetOut() {
		go func(group string, chanOut chan<- container.Container) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					data := p.fsaver.ReadData(ctx, group)
					chanOut <- data
				}
			}
		}(i, j)
	}
}
