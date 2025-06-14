package domain

import (
	"time"

	"github.com/google/uuid"
)

type Priority int

const (
	High Priority = iota
	Medium
	Low
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
)

type Task struct {
	ID        string    `json:"id"`
	Priority  Priority  `json:"priority"`
	Payload   []byte    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
	Status    Status    `json:"status"`
}

func NewTask(priority Priority, payload []byte) *Task {
	return &Task{
		ID:        uuid.NewString(),
		Priority:  priority,
		Payload:   payload,
		CreatedAt: time.Now(),
		Status:    StatusPending,
	}
}
