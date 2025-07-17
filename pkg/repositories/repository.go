package repositories

import (
	"15_07_2025/internal/model"
	"sync"
)

type Repository struct {
	tasks map[string]*model.Task
	mu    sync.RWMutex
}

func NewRepository() *Repository {
	return &Repository{
		tasks: make(map[string]*model.Task),
	}
}
func (s *Repository) SaveTask(task *model.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.ID] = task
}

func (s *Repository) GetTask(id string) (*model.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, ok := s.tasks[id]
	if !ok {
		return nil, model.ErrTaskNotFound
	}
	return task, nil
}

func (s *Repository) GetAllTasks() []*model.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	all := make([]*model.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		all = append(all, t)
	}
	return all
}

func (s *Repository) DeleteTask(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok {
		return model.ErrTaskNotFound
	}
	delete(s.tasks, id)
	return nil
}
