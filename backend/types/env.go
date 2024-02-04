package types

import (
    "github.com/jackc/pgx/v5"
    "task-service/repository"
)

// Database connection should be constructed in main
var Conn *pgx.Conn
// Repository should be constructed in main
var Repository *repository.Queries

