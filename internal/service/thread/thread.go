package thread

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"os/signal"
	"simple_brocker/internal/config"
	"simple_brocker/internal/service/event"
	"simple_brocker/internal/service/most"
	"simple_brocker/internal/service/queue"
	"syscall"
	"time"
)

type (
	thread struct {
		cfg         config.Config
		threadQueue map[string]queue.Queue
	}

	Thread interface {
		AddThread()
		Run(ctx context.Context, ioCh most.Most)
		Close()

		TRun(ioCh most.Most)
	}
)

func New(cfg config.Config) Thread {
	return &thread{
		cfg:         cfg,
		threadQueue: make(map[string]queue.Queue),
	}
}

func (t *thread) AddThread() {
	ms := t.cfg.GetServices()
	for i, j := range ms {
		t.threadQueue[i] = queue.New(&j)
	}
}

func (t *thread) Close() {
	for _, j := range t.threadQueue {
		j.Close()
	}
}

func (t *thread) Run(ctx context.Context, ioCh most.Most) {
	for i, j := range t.threadQueue {
		go j.Producer(ctx, ioCh.GetInChan(i))
		go j.Logging(ctx)
		go j.Consumer(ctx, ioCh.GetOutChan(i))
	}

}

// ---TESTing FUNC---
func (t *thread) TRun(ioCh most.Most) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i, j := range t.threadQueue {
		go j.Producer(ctx, ioCh.GetInChan(i))
		go j.Logging(ctx)
		go j.Consumer(ctx, ioCh.GetOutChan(i))
	}

	go GenEv(ctx, ioCh.GetIn())
	go PrintEv(ctx, ioCh.GetOut())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("Context canceled, stopping server...")
	case <-quit:
		log.Println("Received termination signal, stopping server...")
	}
	t.Close()
}

// ---TESTing FUNC---
func GenEv(ctx context.Context, ch map[string]chan event.Event) {
	for i := 1; i <= 1000; i++ {
		s, _ := json.Marshal(fmt.Sprintf("hui %d", i))
		ev := event.New(s)
		k := fmt.Sprintf("serv%d", rand.Int()%3+1)
		select {
		case <-ctx.Done():
			for _, j := range ch {
				close(j)
			}
			return
		case ch[k] <- ev:
		}
		time.Sleep(5 * time.Millisecond)
	}
	for _, j := range ch {
		close(j)
	}
}

// ---TESTing FUNC---
func PrintEv(ctx context.Context, ch map[string]chan []event.Event) {
	for {
		for i, j := range ch {
			select {
			case <-ctx.Done():
				return
			default:
				select {
				case <-ctx.Done():
					return
				case k := <-j:
					fmt.Println("---Pack ", i, "---")
					for l, p := range k {
						var text string
						json.Unmarshal(p.GetData(), &text)
						fmt.Println(l, "OUT: ", text, " - ", p.GetCreateDateF())
					}
					fmt.Println("---End ", i, "---")
				}
			}
		}
	}
}
