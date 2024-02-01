package types

import (
    "time"
    "github.com/google/uuid"
)

type STATUS string

const (
    RUNNING STATUS = "running"
    STOPPED STATUS = "stopped"
    FAILED STATUS = "failed"
    COMPLETED STATUS = "completed"
)

type Task struct {
	Id     uuid.UUID `json:"id"`
	Status STATUS `json:"status"`
	CompletedAt time.Time `json:"completedAt"`
}

