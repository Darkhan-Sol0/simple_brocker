package batcher

import (
	"bytes"
	"context"
	"simple_brocker/internal/config"
	"simple_brocker/internal/service/container"
	"time"
)

type (
	fsaver struct {
		cfg    config.Config
		tempCh map[string]chan []byte
	}

	Fsaver interface {
		LogData(data container.Container)
		ReadData(ctx context.Context, group string) []byte
	}
)

func New(cfg config.Config) Fsaver {
	tempCh := make(map[string]chan []byte)
	for i := range cfg.GetGroups() {
		tempCh[i] = make(chan []byte, 100)
	}
	return &fsaver{
		cfg:    cfg,
		tempCh: tempCh,
	}
}

func (l *fsaver) LogData(data container.Container) {
	l.tempCh[data.Group] <- data.Data
}

func (l *fsaver) ReadData(ctx context.Context, group string) []byte {
	data := make([][]byte, 0)
	timer := time.NewTimer(l.cfg.GetGroup(group).GetCoolDown())
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			if len(data) > 0 {
				return marshalBytes(data)
			}
			return []byte("[]")
		case d := <-l.tempCh[group]:
			data = append(data, d)
			if len(data) >= l.cfg.GetGroup(group).GetServiceBatchSize() {
				return marshalBytes(data)
			}
			timer.Reset(l.cfg.GetGroup(group).GetCoolDown())
		case <-timer.C:
			if len(data) > 0 {
				return marshalBytes(data)
			}
			timer.Reset(l.cfg.GetGroup(group).GetCoolDown())
		}
	}

}

func marshalBytes(messages [][]byte) []byte {
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
