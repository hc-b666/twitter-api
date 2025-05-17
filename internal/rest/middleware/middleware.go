package middleware

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Middleware struct {
	db *pgxpool.Pool
}

func New(
	db *pgxpool.Pool,
) *Middleware {
	return &Middleware{
		db: db,
	}
}
