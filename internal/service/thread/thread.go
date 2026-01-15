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
		ioCh        most.Most
		threadQueue queue.Queue
	}

	Thread interface {
		Run(ctx context.Context)
		Close()

		TRun(ioCh most.Most)
	}
)

func New(cfg config.Config, ioCh most.Most) Thread {
	return &thread{
		cfg:         cfg,
		ioCh:        ioCh,
		threadQueue: queue.New(cfg),
	}
}

func (t *thread) Close() {
	t.threadQueue.Close()
}

func (t *thread) Run(ctx context.Context) {
	go t.threadQueue.Producer(ctx, t.ioCh.GetIn())
	go t.threadQueue.Logging(ctx)
	go t.threadQueue.Consumer(ctx, t.ioCh.GetOut())
}

// ---TESTing FUNC---
func (t *thread) TRun(ioCh most.Most) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go t.threadQueue.Producer(ctx, ioCh.GetIn())
	go t.threadQueue.Logging(ctx)
	go t.threadQueue.Consumer(ctx, ioCh.GetOut())

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
func GenEv(ctx context.Context, ch chan event.Event) {
	for i := 1; i <= 1000; i++ {
		s, _ := json.Marshal(fmt.Sprintf("hui %d", i))
		ev := event.New(fmt.Sprintf("group%d", rand.Int()%3+1), s)
		select {
		case <-ctx.Done():
			close(ch)
			return
		case ch <- ev:
		}
		time.Sleep(5 * time.Millisecond)
	}

	close(ch)

}

// ---TESTing FUNC---
func PrintEv(ctx context.Context, ch chan []event.Event) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			select {
			case <-ctx.Done():
				return
			case k := <-ch:
				fmt.Println("---Pack ", k[0].GetGroup(), "---")
				for l, p := range k {
					var text string
					json.Unmarshal(p.GetData(), &text)
					fmt.Println(l, "OUT: ", text, " - ", p.GetCreateDateF())
				}
				fmt.Println("---End ", k[0].GetGroup(), "---")
			}
		}
	}
}
