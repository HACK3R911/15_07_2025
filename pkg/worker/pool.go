package worker

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"15_07_2025/internal/model"
	"15_07_2025/pkg/repositories"
)

type WorkerPool struct {
	tasks      chan *model.Task
	storage    *repositories.Repository
	wg         sync.WaitGroup
	maxWorkers int
	active     int
	mu         sync.Mutex
}

func NewWorkerPool(n int, s *repositories.Repository) *WorkerPool {
	return &WorkerPool{
		tasks:      make(chan *model.Task, n),
		storage:    s,
		maxWorkers: n,
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.maxWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.tasks)
	wp.wg.Wait()
}

func (wp *WorkerPool) AddTask(task *model.Task) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	if wp.active < wp.maxWorkers {
		wp.tasks <- task
		wp.active++
	}
}

func (wp *WorkerPool) ActiveTasks() int {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	return wp.active
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	for task := range wp.tasks {
		wp.processTask(task)
		wp.mu.Lock()
		wp.active--
		wp.mu.Unlock()
	}
}

func (wp *WorkerPool) processTask(task *model.Task) {
	archivePath, err := wp.createArchive(task)
	if err != nil {
		task.Status = "failed"
		task.Error = err.Error()
	} else {
		task.Status = "completed"
		task.ArchivePath = archivePath
	}
	wp.storage.SaveTask(task)
}

func (wp *WorkerPool) createArchive(task *model.Task) (string, error) {
	archivePath := filepath.Join(os.TempDir(), task.ID+".zip")
	f, err := os.Create(archivePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	zipWriter := zip.NewWriter(f)
	defer zipWriter.Close()

	var errors []string

	for _, url := range task.URLs {
		if err := wp.addToZip(zipWriter, url); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", url, err))
		}
	}

	if len(errors) > 0 {
		task.Error = "Some downloads failed: " + strings.Join(errors, "; ")
	}

	return archivePath, nil
}

func (wp *WorkerPool) addToZip(w *zip.Writer, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP status %d", resp.StatusCode)
	}

	entry, err := w.Create(filepath.Base(url))
	if err != nil {
		return err
	}

	_, err = io.Copy(entry, resp.Body)
	return err
}
