package repository

import (
	"context"
	"database/sql"
	"fmt"
	"gitlab.com/spacewalker/locations/internal/app/user/core/domain"
	"gitlab.com/spacewalker/locations/internal/app/user/core/port"
)

var createUserQuery = fmt.Sprintf(
	`
INSERT INTO %s
(username)
VALUES ($1)
RETURNING id, username, created_at, updated_at
`,
	UserTable,
)

// CreateUser add new User to database and returns it.
func (q *postgresQueries) CreateUser(ctx context.Context, arg port.CreateUserArg) (domain.User, error) {
	var user domain.User

	if err := q.db.QueryRowContext(ctx, createUserQuery, arg.Username).Scan(
		&user.ID,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

var getByUsernameQuery = fmt.Sprintf(
	`
SELECT id, username, created_at, updated_at
FROM %s
WHERE username = $1
`,
	UserTable,
)

// GetByUsername returns user with given username.
// If user is not found returns ErrNotFound error.
func (q *postgresQueries) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	var user domain.User
	if err := q.db.QueryRowContext(ctx, getByUsernameQuery, username).Scan(
		&user.ID,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, port.ErrNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}

// SetUserLocation gets User by given username and updates Location by user ID
//with provided coordinates within a single database transaction.
func (r *postgresRepository) SetUserLocation(ctx context.Context, arg port.SetUserLocationArg) (domain.Location, error) {
	var location domain.Location

	err := r.execTx(ctx, func(q *postgresQueries) error {
		user, err := q.GetByUsername(ctx, arg.Username)
		if err == port.ErrNotFound {
			user, err = q.CreateUser(ctx, port.CreateUserArg{Username: arg.Username})
			if err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}

		location, err = q.SetLocation(ctx, port.SetLocationArg{
			UserID:    user.ID,
			Latitude:  arg.Latitude,
			Longitude: arg.Longitude,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return domain.Location{}, err
	}

	return location, nil
}
