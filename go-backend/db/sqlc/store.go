package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	connPool *pgxpool.Pool
}

// Creates a new store object
func NewStore(coonPool *pgxpool.Pool) *Store {
	return &Store{
		Queries:  New(coonPool),
		connPool: coonPool,
	}
}
