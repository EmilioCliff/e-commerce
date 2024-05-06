package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
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

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("err: %w rbErr: %w", err, rbErr)
		}
		return fmt.Errorf("err: %w", err)
	}

	return tx.Commit(ctx)
}
