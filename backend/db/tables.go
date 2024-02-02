package db

import (
	"time"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const DB_SCHEMA = `

DO $$
BEGIN 
    IF NOT EXISTS (select 1 from pg_type where typname = 'task_status')  THEN
        CREATE TYPE task_status AS ENUM (
            'queued',
            'running',
            'failed'
            'completed'
        );
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS task (
    id uuid PRIMARY KEY,
    status task_status NOT NULL,
    completed_at timestamp
);`

type Task struct {
	Id          uuid.UUID `db:"id"`
	Status      string    `db:"status"`
	CompletedAt time.Time `db:"completed_at"`
}
