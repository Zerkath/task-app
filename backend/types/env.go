package types

import (
    "task-service/repository"
    "github.com/jackc/pgx/v5/pgxpool"
)

// Database connection should be constructed in main
var Conn *pgxpool.Pool
// Repository should be constructed in main
var Repository *repository.Queries

