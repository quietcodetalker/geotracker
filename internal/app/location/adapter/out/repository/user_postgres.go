package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"gitlab.com/spacewalker/locations/internal/app/location/core/domain"
	"gitlab.com/spacewalker/locations/internal/app/location/core/port"
	"gitlab.com/spacewalker/locations/internal/pkg/geo"
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
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code.Name() {
			case "string_data_right_truncation":
				return domain.User{}, port.ErrInvalidUsername
			}

			switch pqErr.Constraint {
			case ConstraintUsersUsernameKey:
				return domain.User{}, port.ErrAlreadyExists
			case ConstraintUsersUsernameValid:
				return domain.User{}, port.ErrInvalidUsername
			}
		}
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
func (r *postgresRepository) SetUserLocation(ctx context.Context, arg port.UserRepositorySetUserLocationRequest) (port.UserRepositorySetUserLocationResponse, error) {
	var user domain.User
	var prevLocation domain.Location
	var location domain.Location

	err := r.execTx(ctx, func(q *postgresQueries) error {
		var err error

		user, err = q.GetByUsername(ctx, arg.Username)
		if err == nil {
			var glErr error
			prevLocation, glErr = q.GetLocation(ctx, user.ID)
			if glErr != nil && glErr != port.ErrNotFound {
				return err
			}
		}
		if err == port.ErrNotFound {
			user, err = q.CreateUser(ctx, port.CreateUserArg{Username: arg.Username})
			if err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}

		location, err = q.SetLocation(ctx, port.LocationRepositorySetLocationRequest{
			UserID: user.ID,
			Point:  arg.Point,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return port.UserRepositorySetUserLocationResponse{}, err
	}

	return port.UserRepositorySetUserLocationResponse{
		User:         user,
		PrevLocation: prevLocation,
		Location:     location,
	}, nil
}

var listUsersInRadiusQuery = fmt.Sprintf(
	`
SELECT u.id, u.username, u.created_at, u.updated_at
FROM %s u
INNER JOIN %s l ON l.user_id = u.id
WHERE ($1<@>l.point) * 1609.344 <= $2 AND u.id > $3
LIMIT $4
`,
	UserTable,
	LocationTable,
)

// ListUsersInRadius retrieve users in given radius with coordinates.
func (q *postgresQueries) ListUsersInRadius(ctx context.Context, arg port.UserRepositoryListUsersInRadiusRequest) (port.UserRepositoryListUsersInRadiusResponse, error) {
	var users []domain.User

	// Fetch PageSize + 1 (extra marker element)
	// If such element happens to be retrieved it means that next page can be (probably) retrieved as well.
	rows, err := q.db.QueryContext(ctx, listUsersInRadiusQuery, geo.PostgresPoint(arg.Point), arg.Radius, arg.PageToken, arg.PageSize+1)
	if err != nil {
		if err == sql.ErrNoRows {
			return port.UserRepositoryListUsersInRadiusResponse{}, port.ErrNotFound
		}
		return port.UserRepositoryListUsersInRadiusResponse{}, err
	}
	defer rows.Close()

	hasNextPage := false
	counter := 0
	for rows.Next() {
		counter++
		if counter > arg.PageSize { // Next page exists.
			hasNextPage = true
			break // Do not scan extra marker element.
		}

		var user domain.User
		if err = rows.Scan(
			&user.ID,
			&user.Username,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return port.UserRepositoryListUsersInRadiusResponse{}, err
		}
		users = append(users, user)
	}

	result := port.UserRepositoryListUsersInRadiusResponse{
		Users: users,
	}
	if hasNextPage {
		result.NextPageToken = users[len(users)-1].ID
	}

	return result, nil
}
