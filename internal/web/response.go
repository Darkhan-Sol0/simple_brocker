package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"simple_brocker/internal/service/event"
	"time"
)

func (r *router) ResponseEvent(ctx context.Context, ch chan []event.Event) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case temp := <-ch:
			buff := make([]json.RawMessage, 0)

			for _, j := range temp {
				buff = append(buff, j.GetData())
			}
			payload := marshalRawMessages(buff)

			grcfg := r.cfg.GetGroup(temp[0].GetGroup())

			for _, j := range grcfg.GetServiceAddress() {
				go func(address string, data []byte) {
					for attempt := 1; attempt <= grcfg.GetRetry(); attempt++ {
						err := senderMessage(address, data)
						if err == nil {
							return
						}
						log.Printf("Attempt %d/%d failed for %s: %v",
							attempt, grcfg.GetRetry(), address, err)
						if attempt < grcfg.GetRetry() {
							backoff := time.Duration(attempt*attempt) * 100 * time.Millisecond
							time.Sleep(backoff)
						}
					}
					log.Printf("All %d attempts failed for %s", grcfg.GetRetry(), address)
				}(j, payload)
			}
			fmt.Println(buff)
			var tes any
			json.Unmarshal(payload, &tes)
			fmt.Println(tes)
		}
	}
}

func senderMessage(address string, data []byte) error {
	req, _ := http.NewRequest("POST", address, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Request to %s failed: %v", address, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Bad response from %s: %d - %s", address, resp.StatusCode, string(body))
	}
	return nil
}

func marshalRawMessages(messages []json.RawMessage) []byte {
	if len(messages) == 0 {
		return []byte("[]")
	}
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	buf.WriteByte('[')
	buf.Write(messages[0])
	for i := 1; i < len(messages); i++ {
		buf.WriteByte(',')
		buf.Write(messages[i])
	}
	buf.WriteByte(']')
	return buf.Bytes()
}
