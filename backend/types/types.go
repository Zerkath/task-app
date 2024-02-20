package types

import (
    "github.com/google/uuid"
    "task-service/repository"
)

type STATUS string

const (
    QUEUED STATUS = "queued"
    RUNNING STATUS = "running"
    FAILED STATUS = "failed"
    COMPLETED STATUS = "completed"
)

type Task struct {
	Id     uuid.UUID `json:"id"`
	Status STATUS `json:"status"`
    // should be omitted if nil
	CompletedAt int32 `json:"completedAt,omitempty"`
}

type Page struct {
    Count int64 `json:"count"`
    Page int `json:"page"`
    Data []repository.Task `json:"data"`
}
