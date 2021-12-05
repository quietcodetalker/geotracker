package repository_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"gitlab.com/spacewalker/locations/internal/app/user/adapter/out/repository"
	"gitlab.com/spacewalker/locations/internal/app/user/core/domain"
	"gitlab.com/spacewalker/locations/internal/app/user/core/port"
)

var seedUsersQuery = fmt.Sprintf(
	`
INSERT INTO %s
(username)
VALUES ($1)
RETURNING id, username, created_at, updated_at
`,
	repository.UserTable,
)

func (s *PostgresTestSuite) seedUsers(args []port.CreateUserArg) []domain.User {
	users := make([]domain.User, 0, len(args))

	tx, err := s.db.Begin()
	require.NoError(s.T(), err)
	defer tx.Rollback()

	stmt, err := tx.Prepare(seedUsersQuery)
	require.NoError(s.T(), err)
	defer stmt.Close()

	for _, arg := range args {
		var user domain.User

		err = stmt.QueryRow(arg.Username).Scan(
			&user.ID,
			&user.Username,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		require.NoError(s.T(), err)

		users = append(users, user)
	}

	err = tx.Commit()
	require.NoError(s.T(), err)

	return users
}

var seedLocationsQuery = fmt.Sprintf(
	`
INSERT INTO %s
(user_id, latitude, longitude)
VALUES ($1, $2, $3)
RETURNING user_id, latitude, longitude, created_at, updated_at
`,
	repository.LocationTable,
)

func (s *PostgresTestSuite) seedLocations(args []port.SetLocationArg) []domain.Location {
	locations := make([]domain.Location, 0, len(args))

	tx, err := s.db.Begin()
	require.NoError(s.T(), err)

	defer tx.Rollback()

	stmt, err := tx.Prepare(seedLocationsQuery)
	require.NoError(s.T(), err)
	defer stmt.Close()

	for _, arg := range args {
		var location domain.Location

		err := stmt.QueryRow(arg.UserID, arg.Latitude, arg.Longitude).Scan(
			&location.UserID,
			&location.Latitude,
			&location.Longitude,
			&location.CreatedAt,
			&location.UpdatedAt,
		)
		require.NoError(s.T(), err)

		locations = append(locations, location)
	}

	err = tx.Commit()
	require.NoError(s.T(), err)

	return locations
}
