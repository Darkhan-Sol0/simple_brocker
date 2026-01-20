package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"simple_brocker/internal/service/event"
	"sync"
	"time"
)

type (
	lnifiedLogger struct {
		baseDir    string
		files      map[string]*os.File
		currentDay string
		mu         sync.RWMutex
	}

	Logger interface {
		LogEvent(ev event.EventIn) error
		Close() error
	}

	LogEntry struct {
		Time  string          `json:"time"`
		Group string          `json:"group"`
		Data  json.RawMessage `json:"data"`
	}
)

func New() (Logger, error) {
	if err := os.MkdirAll("./logs/", 0755); err != nil {
		return nil, fmt.Errorf("failed to create base log directory: %w", err)
	}

	return &lnifiedLogger{
		baseDir:    "./logs/",
		files:      make(map[string]*os.File),
		currentDay: time.Now().Format("2006-01-02"),
	}, nil
}

func (ul *lnifiedLogger) ensureGroupFile(group string) (*os.File, error) {
	ul.mu.Lock()
	defer ul.mu.Unlock()
	today := time.Now().Format("2006-01-02")
	if ul.currentDay != today {
		for groupName, file := range ul.files {
			file.Close()
			delete(ul.files, groupName)
		}
		ul.currentDay = today
	}
	if file, exists := ul.files[group]; exists {
		return file, nil
	}
	groupDir := filepath.Join(ul.baseDir, group)
	if err := os.MkdirAll(groupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create group directory: %w", err)
	}
	filePath := filepath.Join(groupDir, ul.currentDay+".log")
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file for group %s: %w", group, err)
	}

	ul.files[group] = file
	return file, nil
}

func (ul *lnifiedLogger) LogEvent(ev event.EventIn) error {
	file, err := ul.ensureGroupFile(ev.GetGroup())
	if err != nil {
		return err
	}
	ul.mu.RLock()
	defer ul.mu.RUnlock()

	entry := LogEntry{
		Time:  time.Now().Format("15:04:05.000000"),
		Group: ev.GetGroup(),
		Data:  ev.GetData(),
	}
	jsonData, _ := json.Marshal(entry)
	jsonData = append(jsonData, '\n')
	_, err = file.WriteString(string(jsonData))
	return err
}

func (ul *lnifiedLogger) Close() error {
	ul.mu.Lock()
	defer ul.mu.Unlock()
	var lastErr error
	for group, file := range ul.files {
		if err := file.Close(); err != nil {
			lastErr = err
			fmt.Printf("Error closing log file for group %s: %v\n", group, err)
		}
		delete(ul.files, group)
	}
	return lastErr
}
