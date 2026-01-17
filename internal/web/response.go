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

func (r *router) ResponseEvent(ctx context.Context, ch <-chan []event.Event) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case temp := <-ch:
			buff := make([]json.RawMessage, 0)

			for _, j := range temp {
				buff = append(buff, j.GetData())
			}

			for _, j := range r.cfg.GetGroup(temp[0].GetGroup()).GetServiceAddress() {
				go func(address string, data []json.RawMessage) {
					payload, _ := json.Marshal(data)

					req, _ := http.NewRequest("POST", address, bytes.NewBuffer(payload))

					req.Header.Set("Content-Type", "application/json")

					client := &http.Client{
						Timeout: 5 * time.Second,
					}

					resp, err := client.Do(req)
					if err != nil {
						log.Printf("Request to %s failed: %v", address, err)
						return
					}
					defer resp.Body.Close()
					if resp.StatusCode >= 400 {
						body, _ := io.ReadAll(resp.Body)
						log.Printf("Bad response from %s: %d - %s", address, resp.StatusCode, string(body))
					}

				}(j, buff)
			}

			// type aaaa struct {
			// 	Data string `json:"data"`
			// 	Num  int    `json:"num"`
			// 	Bool bool   `json:"bool"`
			// }

			payload, _ := json.Marshal(buff)

			fmt.Println(buff)
			var tes any
			json.Unmarshal(payload, &tes)
			fmt.Println(tes)
		}
	}
}
