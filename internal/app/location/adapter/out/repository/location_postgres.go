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

var setLocationQuery = fmt.Sprintf(
	`
INSERT INTO %s
(user_id, point)
VALUES ($1, $2)
ON CONFLICT ON CONSTRAINT locations_pkey 
DO
	UPDATE SET point = EXCLUDED.point
RETURNING user_id, point, created_at, updated_at
`,
	LocationTable,
)

// SetLocation adds a new record to locations table with `arg.UserID` as `user_id` and
// `arg.Point` as `point`.
// If a record with given id already exists it only updates point.
//
// Returns `Location` populated with data from the created or updated record and `error`.
//
// `ErrFailedPrecondition` is returned, in case there is no user with given id.
//
// `ErrInvalidArgument` is returned in case given point's longitude or latitude is invalid.,
//
// `ErrInvalidError` is returned in case of any other failure.
//
// Returned error is wrapped with `fmt.Errorf("%w", err)`, so use `errors.Is()` to compare returned error.
func (q *postgresQueries) SetLocation(ctx context.Context, arg port.LocationRepositorySetLocationRequest) (domain.Location, error) {
	var location domain.Location
	var pgPoint geo.PostgresPoint
	if err := q.db.QueryRowContext(ctx, setLocationQuery, arg.UserID, geo.PostgresPoint(arg.Point)).Scan(
		&location.UserID,
		&pgPoint,
		&location.CreatedAt,
		&location.UpdatedAt,
	); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Constraint {
			case ConstraintLocationsUserIdFkey:
				return domain.Location{}, fmt.Errorf("%w", errpack.ErrFailedPrecondition)
			case ConstraintLocationsLatitudeValid:
				return domain.Location{}, fmt.Errorf("%w", errpack.ErrInvalidArgument)
			case ConstraintLocationsLongitudeValid:
				return domain.Location{}, fmt.Errorf("%w", errpack.ErrInvalidArgument)
			}
		}
		return domain.Location{}, fmt.Errorf("%w: %v", errpack.ErrInternalError, err)
	}

	location.Point = geo.Point(pgPoint)

	return location, nil
}

var getLocationQuery = fmt.Sprintf(
	`
SELECT user_id, point, created_at, updated_at
FROM %s
WHERE user_id = $1
`,
	LocationTable,
)

// GetLocation finds a location by given user id in the locations table.
//
// It returns a found location and any error encountered.
//
// `ErrNotFound` is returned in case required location is not found.
//
// `ErrInternalError` is returned in case of any other error.
//
// Returned error is wrapped with `fmt.Errorf("%w", err)`, so use `errors.Is()` to compare returned error.
func (q *postgresQueries) GetLocation(ctx context.Context, userID int) (domain.Location, error) {
	var location domain.Location
	var point geo.PostgresPoint

	if err := q.db.QueryRowContext(ctx, getLocationQuery, userID).Scan(
		&location.UserID,
		&point,
		&location.CreatedAt,
		&location.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Location{}, fmt.Errorf("%w", errpack.ErrNotFound)
		}
		return domain.Location{}, fmt.Errorf("%w: %v", errpack.ErrInternalError, err)
	}

	location.Point = geo.Point(point)

	return location, nil
}
