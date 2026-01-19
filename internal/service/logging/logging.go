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
	// lnifiedLogger - один логгер для всех групп
	lnifiedLogger struct {
		baseDir    string
		files      map[string]*os.File // группа -> файл
		currentDay string              // текущий день "02-01-2006"
		mu         sync.RWMutex
	}

	Logger interface {
		LogEvent(ev event.Event) error
		Close() error
	}

	LogEntry struct {
		Time  string          `json:"time"`
		Group string          `json:"group"`
		Data  json.RawMessage `json:"data"`
	}
)

// New создает единый логгер для всех групп
func New() (Logger, error) {
	if err := os.MkdirAll("./logs/", 0755); err != nil {
		return nil, fmt.Errorf("failed to create base log directory: %w", err)
	}

	return &lnifiedLogger{
		baseDir:    "./logs/",
		files:      make(map[string]*os.File),
		currentDay: time.Now().Format("02-01-2006"),
	}, nil
}

// ensureGroupFile проверяет и открывает файл для конкретной группы
func (ul *lnifiedLogger) ensureGroupFile(group string) (*os.File, error) {
	ul.mu.Lock()
	defer ul.mu.Unlock()

	// Проверяем не сменился ли день
	today := time.Now().Format("2006-01-02")
	if ul.currentDay != today {
		// День сменился - закрываем все старые файлы
		for groupName, file := range ul.files {
			file.Close()
			delete(ul.files, groupName)
		}
		ul.currentDay = today
	}

	// Если файл для этой группы уже открыт - возвращаем
	if file, exists := ul.files[group]; exists {
		return file, nil
	}

	// Создаем поддиректорию для группы
	groupDir := filepath.Join(ul.baseDir, group)
	if err := os.MkdirAll(groupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create group directory: %w", err)
	}

	// Открываем файл для группы
	filePath := filepath.Join(groupDir, ul.currentDay+".log")
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file for group %s: %w", group, err)
	}

	ul.files[group] = file
	return file, nil
}

func (ul *lnifiedLogger) LogEvent(ev event.Event) error {
	// Получаем или создаем файл для этой группы
	file, err := ul.ensureGroupFile(ev.GetGroup())
	if err != nil {
		return err
	}

	ul.mu.RLock()
	defer ul.mu.RUnlock()

	entry := LogEntry{
		Time:  time.Now().Format("15:04:05.000000"),
		Group: ev.GetGroup(),
		Data:  ev.GetData(), // Оригинальный JSON как есть
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
