package model

import (
	"errors"
	"strconv"
	"time"
)

var (
	ErrTaskNotFound    = errors.New("task not found")
	ErrServerBusy      = errors.New("server is busy, try again later")
	ErrMaxURLsReached  = errors.New("maximum number of URLs (3) reached")
	ErrArchiveNotReady = errors.New("archive is not ready yet")
	ErrInvalidFileType = errors.New("invalid file type")
)

type Task struct {
	ID          string    `json:"id"`
	URLs        []string  `json:"urls"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	ArchivePath string    `json:"archive_path,omitempty"`
	Error       string    `json:"error,omitempty"`
}

func NewTask() *Task {
	return &Task{
		ID:        generateTaskID(),
		Status:    "pending",
		CreatedAt: time.Now(),
	}
}

func (t *Task) AddURL(url string) {
	t.URLs = append(t.URLs, url)
}

func generateTaskID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
