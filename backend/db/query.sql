-- name: GetTaskById :one
SELECT 
    id, status, created_at, completed_at, restarts
FROM task
WHERE id = $1 LIMIT 1;

-- name: GetTaskCount :one
SELECT COUNT(*)
FROM task
WHERE sqlc.narg('status')::task_status IS NULL
OR status in (sqlc.narg('status')::task_status);

-- name: GetTasks :many
SELECT
    id, status, created_at, completed_at, restarts
FROM task
WHERE sqlc.narg('status')::task_status IS NULL
OR status in (sqlc.narg('status')::task_status)
ORDER BY status, created_at, id DESC
LIMIT $1 OFFSET $2;

-- name: GetTasksByIds :many
SELECT
    id, status, created_at, completed_at, restarts
FROM task
WHERE id = ANY(sqlc.arg('ids')::uuid[])
ORDER BY status, created_at, id DESC;

-- name: GetUncompletedTasks :many
SELECT
    id, status, created_at, completed_at, restarts
FROM task
WHERE status != 'completed' 
ORDER BY status, created_at, id DESC
LIMIT $1 OFFSET $2;

-- name: NewTask :one
INSERT INTO task
DEFAULT VALUES
RETURNING id, status, created_at, completed_at, restarts;

-- name: UpdateTask :one
UPDATE task
SET status = $1, completed_at = $2, restarts = $3
WHERE id = $4
RETURNING id, status, created_at, completed_at, restarts;

-- name: DeleteTask :exec
DELETE FROM task
WHERE id = $1;

