package handler

import (
	"encoding/json"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"15_07_2025/config"
	"15_07_2025/internal/model"
	"15_07_2025/pkg/repositories"
	"15_07_2025/pkg/worker"
)

type Handler struct {
	repositories *repositories.Repository
	workerPool   *worker.WorkerPool
	cfg          *config.Config
}

func NewArchiveHandler(s *repositories.Repository, w *worker.WorkerPool, cfg *config.Config) *Handler {
	return &Handler{s, w, cfg}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/tasks":
		h.createTask(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/tasks":
		h.listTasks(w, r)
	case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/tasks/"):
		h.deleteTask(w, r)
	case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/urls"):
		h.addURL(w, r)
	case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/download"):
		h.download(w, r)
	case r.Method == http.MethodGet:
		h.getStatus(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	if h.workerPool.ActiveTasks() >= 3 {
		http.Error(w, model.ErrServerBusy.Error(), http.StatusServiceUnavailable)
		return
	}

	task := model.NewTask()
	h.repositories.SaveTask(task)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)

}

func (h *Handler) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.repositories.GetAllTasks()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.Path)

	err := h.repositories.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) addURL(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.NotFound(w, r)
		return
	}
	taskID := pathParts[2]

	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ext := strings.ToLower(filepath.Ext(req.URL))
	if ext == "" {
		http.Error(w, model.ErrInvalidFileType.Error(), http.StatusBadRequest)
		return
	}
	ext = ext[1:]

	if !h.cfg.AllowedExts[ext] {
		http.Error(w, model.ErrInvalidFileType.Error(), http.StatusBadRequest)
		return
	}

	task, err := h.repositories.GetTask(taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if len(task.URLs) >= 3 {
		http.Error(w, model.ErrMaxURLsReached.Error(), http.StatusBadRequest)
		return
	}

	task.AddURL(req.URL)
	h.repositories.SaveTask(task)

	if len(task.URLs) == 3 {
		h.workerPool.AddTask(task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) getStatus(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.NotFound(w, r)
		return
	}
	taskID := pathParts[2]

	task, err := h.repositories.GetTask(taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) download(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.NotFound(w, r)
		return
	}
	taskID := pathParts[2]

	task, err := h.repositories.GetTask(taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if task.ArchivePath == "" {
		http.Error(w, model.ErrArchiveNotReady.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(task.ArchivePath))
	w.Header().Set("Content-Type", "application/zip")
	http.ServeFile(w, r, task.ArchivePath)
}
