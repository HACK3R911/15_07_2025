package repositories

import (
	"15_07_2025/internal/model"
	"sync"
)

type MemoryStorage struct {
	tasks map[string]*model.Task
	mu    sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tasks: make(map[string]*model.Task),
	}
}
