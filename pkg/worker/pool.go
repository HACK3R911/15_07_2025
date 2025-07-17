package worker

import (
	"15_07_2025/internal/model"
	"15_07_2025/pkg/repositories"
	"sync"
)

type WorkerPool struct {
	tasks      chan *model.Task
	storage    *repositories.MemoryStorage
	wg         sync.WaitGroup
	maxWorkers int
	active     int
	mu         sync.Mutex
}

func NewWorkerPool(n int, s *repositories.MemoryStorage) *WorkerPool {
	return &WorkerPool{
		tasks:      make(chan *model.Task, n),
		storage:    s,
		maxWorkers: n,
	}
}
