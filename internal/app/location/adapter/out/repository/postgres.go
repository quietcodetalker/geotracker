package repository

import (
	"context"
	"database/sql"
	"fmt"
	"gitlab.com/spacewalker/locations/internal/app/location/core/port"
	"gitlab.com/spacewalker/locations/internal/pkg/errpack"
)

const (
	ConstraintUsersUsernameKey   = "users_username_key"
	ConstraintUsersUsernameValid = "users_username_valid"

	ConstraintLocationsUserIdFkey     = "locations_user_id_fkey"
	ConstraintLocationsLatitudeValid  = "locations_latitude_valid"
	ConstraintLocationsLongitudeValid = "locations_longitude_valid"
)

type postgresRepository struct {
	*postgresQueries
	db *sql.DB
}

// NewPostgresRepository returns a new instance of port.Repository
func NewPostgresRepository(db *sql.DB) port.Repository {
	return &postgresRepository{
		postgresQueries: newPostgresQueries(db),
		db:              db,
	}
}

// execTx executes provided callback in the scope of a database transaction.
//
// It returns an error occurred while starting, committing or rolling back the transaction or
// and error returned by the callback.
func (r *postgresRepository) execTx(ctx context.Context, fn func(*postgresQueries) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%w: %v", errpack.ErrInternalError, err)
	}

	q := newPostgresQueries(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("%w: %v", errpack.ErrInternalError, err)
		}
		return fmt.Errorf("%w: %v", errpack.ErrInternalError, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%w: %v", errpack.ErrInternalError, err)
	}

	return nil
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
