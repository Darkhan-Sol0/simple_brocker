package response

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"simple_brocker/internal/config"
	"simple_brocker/internal/service/container"
	"time"
)

type (
	response struct {
		cfg     config.Config
		chanOut map[string]chan container.Container
	}

	Response interface {
		Sender(ctx context.Context)
	}
)

func New(cfg config.Config, chanOut map[string]chan container.Container) Response {
	return &response{
		cfg:     cfg,
		chanOut: chanOut,
	}
}

func (r *response) Sender(ctx context.Context) {
	for i, j := range r.chanOut {
		go func(group string, chanOut <-chan container.Container) {
			for {
				select {
				case <-ctx.Done():
					return
				case data := <-chanOut:
					address := r.cfg.GetGroup(group).GetServiceAddress()
					for _, k := range address {
						go retrySender(ctx, k, r.cfg.GetGroup(group).GetRetry(), data.Data)
					}

				}
			}
		}(i, j)
	}
}

func senderMessage(ctx context.Context, address string, data []byte) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		req, err := http.NewRequestWithContext(ctx, "POST", address, bytes.NewBuffer(data))
		if err != nil {
			return fmt.Errorf("err: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{
			Timeout: 5 * time.Second,
		}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("err: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("Bad response from %s: %d - %s", address, resp.StatusCode, string(body))
		}
		return nil
	}
}

func retrySender(ctx context.Context, address string, try int, data []byte) error {
	for attempt := 1; attempt <= try; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := senderMessage(ctx, address, data)
			if err == nil {
				return nil
			}
			log.Printf("Attempt %d/%d failed for %s: %v",
				attempt, try, address, err)
			if attempt < try {
				backoff := time.Duration(attempt*attempt) * 100 * time.Millisecond
				time.Sleep(backoff)
			}
		}
	}
	log.Printf("All %d attempts failed for %s", try, address)
	return nil
}
