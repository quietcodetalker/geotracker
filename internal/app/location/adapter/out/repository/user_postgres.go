package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"gitlab.com/spacewalker/locations/internal/app/location/core/domain"
	"gitlab.com/spacewalker/locations/internal/app/location/core/port"
	"gitlab.com/spacewalker/locations/internal/pkg/errpack"
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

// CreateUser adds a new user to the users table.
//
// It returns the created user and any error encountered.
//
// `ErrInvalidArgument` is returned in case username is invalid.
//
// `ErrAlreadyExists` is returned in case a user with given username already exists.
//
// `ErrInternalError` is returned in case of any other failure.
//
// Returned error is wrapped with `fmt.Errorf("%w", err)`. Use `errors.Is()` to compare errors.
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
				return domain.User{}, fmt.Errorf("%w", errpack.ErrInvalidArgument)
			}

			switch pqErr.Constraint {
			case ConstraintUsersUsernameKey:
				return domain.User{}, fmt.Errorf("%w", errpack.ErrAlreadyExists)
			case ConstraintUsersUsernameValid:
				return domain.User{}, fmt.Errorf("%w", errpack.ErrInvalidArgument)
			}
		}
		return domain.User{}, fmt.Errorf("%w", errpack.ErrInternalError)
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

// GetByUsername finds a user by username in users table.
//
// It returns `User` and `error`.
//
// `ErrNotFound` is returned in case user not found.
//
// `ErrInternalError` is returned in case of any other failure.
//
// Returned error is wrapped with `fmt.Errorf("%w", err)`. Use `errors.Is()` to compare errors.
func (q *postgresQueries) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	var user domain.User
	if err := q.db.QueryRowContext(ctx, getByUsernameQuery, username).Scan(
		&user.ID,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, fmt.Errorf("%w", errpack.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("%w", errpack.ErrInternalError)
	}

	return user, nil
}

// SetUserLocation sets user's location.
//
// It finds a user by the provided username. If the user is not found, it creates new one.
// If user was found, it finds current location of the user.
// Then sets location of the user. All of it is done in the scope of the database transaction.
//
// It returns a response and any error encountered.
//
// The response consists of:
//		- found or created user
//		- previous location of the user (should be considered as not found if its `UserID` equals 0)
//		- new location of the user
//
// `ErrInternalError` is returned in following cases:
//		- any error encountered while
// 			starting, committing and rolling back the database transaction.
//		- `ErrInternalError` is returned from `GetByUsername`, `ErrInternalError`,
//			`GetLocation` or `CreateUser` methods
//
//	`ErrInvalidArgument` is returned in case `ErrInvalidArgument` is returned from
//	`CreateUser` or `SetLocation` methods.
//
// Returned error is wrapped with `fmt.Errorf("%w", err)`. Use `errors.Is()` to compare errors.
func (r *postgresRepository) SetUserLocation(ctx context.Context, arg port.UserRepositorySetUserLocationRequest) (port.UserRepositorySetUserLocationResponse, error) {
	var user domain.User
	var prevLocation domain.Location
	var location domain.Location

	err := r.execTx(ctx, func(q *postgresQueries) error {
		var err error

		user, err = q.GetByUsername(ctx, arg.Username)
		if err == nil {
			// User is found.
			var glErr error
			prevLocation, glErr = q.GetLocation(ctx, user.ID)
			if glErr != nil && !errors.Is(glErr, errpack.ErrNotFound) {
				return fmt.Errorf("%w", errpack.ErrInternalError)
			}
		}
		if errors.Is(err, errpack.ErrNotFound) {
			// ErrNotFound occurred.
			user, err = q.CreateUser(ctx, port.CreateUserArg{Username: arg.Username})
			if err != nil {
				// ErrInternalError or ErrInvalidArgument occurred.
				return err
			}
		}
		if err != nil {
			// ErrInternalError occurred.
			return fmt.Errorf("%w", errpack.ErrInternalError)
		}

		location, err = q.SetLocation(ctx, port.LocationRepositorySetLocationRequest{
			UserID: user.ID,
			Point:  arg.Point,
		})
		if err != nil {
			// ErrInvalidArgument or ErrInternalErr occurred.
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

// ListUsersInRadius finds no more than `arg.PageSize` users by given radius and coordinates
// with IDs greater than `arg.PageToken` ordered by ID.
//
// It returns a response and any error encountered.
//
// The response consists of a user list and next page token.
// Next page token is ID of last found user if required amount of users found.
//
// A user list that equals nil should be considered as empty.
// Next page id is an ID of the first user of the next page.
// If the next page token equal 0, there is no more pages.
//
// `ErrInternalErr` is returned in case any error encountered.
//
// Returned error is wrapped with `fmt.Errorf("%w", err)`. Use `errors.Is()` to compare errors.
func (q *postgresQueries) ListUsersInRadius(ctx context.Context, arg port.UserRepositoryListUsersInRadiusRequest) (port.UserRepositoryListUsersInRadiusResponse, error) {
	var users []domain.User

	// Fetch PageSize + 1 (extra marker element)
	// If such element happens to be retrieved it means that next page can be (probably) retrieved as well.
	rows, err := q.db.QueryContext(ctx, listUsersInRadiusQuery, geo.PostgresPoint(arg.Point), arg.Radius, arg.PageToken, arg.PageSize+1)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return port.UserRepositoryListUsersInRadiusResponse{}, nil
		}
		return port.UserRepositoryListUsersInRadiusResponse{}, fmt.Errorf("%w", errpack.ErrInternalError)
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
			return port.UserRepositoryListUsersInRadiusResponse{}, fmt.Errorf("%w", errpack.ErrInternalError)
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
