package db

import (
	"time"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const DB_SCHEMA = `
CREATE TYPE IF NOT EXISTS task_status AS ENUM (
    "running",
    "stopped",
    "failed"
    "completed"
);


CREATE TABLE IF NOT EXISTS task (
    id uuid PRIMARY KEY,
    status task_status NOT NULL
    completed_at timestamp
);`

type Task struct {
	Id          uuid.UUID `db:"id"`
	Status      string    `db:"status"`
	CompletedAt time.Time `db:"completed_at"`
}
