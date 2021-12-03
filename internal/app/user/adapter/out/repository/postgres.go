package repository

import (
	"context"
	"database/sql"
	"fmt"
	"gitlab.com/spacewalker/locations/internal/app/user/core/port"
)

type postgresRepository struct {
	*postgresQueries
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) port.Repository {
	return &postgresRepository{
		postgresQueries: newPostgresQueries(db),
		db:              db,
	}
}

func (r *postgresRepository) execTx(ctx context.Context, fn func(*postgresQueries) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := newPostgresQueries(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); err != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type postgresQueries struct {
	db DBTX
}

func (q *postgresQueries) withTx(tx *sql.Tx) *postgresQueries {
	return &postgresQueries{
		db: tx,
	}
}

func newPostgresQueries(db DBTX) *postgresQueries {
	return &postgresQueries{
		db: db,
	}
}
