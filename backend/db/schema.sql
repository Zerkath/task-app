CREATE TYPE task_status AS ENUM (
    'queued',
    'running',
    'failed',
    'completed'
);

CREATE TABLE IF NOT EXISTS task (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    status task_status NOT NULL DEFAULT 'queued',
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    completed_at timestamp DEFAULT NULL,
    restarts int DEFAULT 0
);

