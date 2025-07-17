package handler

import (
	"15_07_2025/config"
	"15_07_2025/pkg/repositories"
	"15_07_2025/pkg/worker"
)

type Handler struct {
	storage *repositories.MemoryStorage
	workers *worker.WorkerPool
	cfg     *config.Config
}

func NewArchiveHandler(s *repositories.MemoryStorage, w *worker.WorkerPool, cfg *config.Config) *Handler {
	return &Handler{s, w, cfg}
}
