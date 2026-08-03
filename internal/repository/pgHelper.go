package repository

import "github.com/jackc/pgx/v5/pgxpool"

type PoolHandler struct {
	pool *pgxpool.Pool
}
