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

// SetLocation sets user's location by given user ID.
// Returns updated location entity.
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
				return domain.Location{}, port.ErrAttemptedSettingLocationOfNonExistentUser
			case ConstraintLocationsLatitudeValid:
				return domain.Location{}, &port.InvalidLocationError{
					Violations: []port.InvalidLocationErrorViolation{
						{
							Subject: "latitude",
							Value:   arg.Point.Latitude(),
						},
					},
				}
			case ConstraintLocationsLongitudeValid:
				return domain.Location{}, &port.InvalidLocationError{
					Violations: []port.InvalidLocationErrorViolation{
						{
							Subject: "longitude",
							Value:   arg.Point.Longitude(),
						},
					},
				}
			}
		}
		return domain.Location{}, err
	}

	location.Point = geo.Point(pgPoint)

	return location, nil
}

var updateLocationyUserIDQuery = fmt.Sprintf(
	`
UPDATE %s
SET point = $2
WHERE user_id = $1
RETURNING user_id, point, created_at, updated_at
`,
	LocationTable,
)

// UpdateLocationByUserID TODO: add description
func (q *postgresQueries) UpdateLocationByUserID(ctx context.Context, arg port.LocationRepositoryUpdateLocationByUserIDRequest) (domain.Location, error) {
	var location domain.Location
	var pgPoint geo.PostgresPoint

	row := q.db.QueryRowContext(ctx, updateLocationyUserIDQuery, arg.UserID, geo.PostgresPoint(arg.Point))

	if err := row.Scan(
		&location.UserID,
		&pgPoint,
		&location.CreatedAt,
		&location.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return domain.Location{}, port.ErrNotFound
		}

		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Constraint {
			case ConstraintLocationsLatitudeValid:
				return domain.Location{}, &port.InvalidLocationError{
					Violations: []port.InvalidLocationErrorViolation{
						{
							Subject: "latitude",
							Value:   arg.Point.Latitude(),
						},
					},
				}
			case ConstraintLocationsLongitudeValid:
				return domain.Location{}, &port.InvalidLocationError{
					Violations: []port.InvalidLocationErrorViolation{
						{
							Subject: "longitude",
							Value:   arg.Point.Longitude(),
						},
					},
				}
			}
		}
		return domain.Location{}, err
	}

	location.Point = geo.Point(pgPoint)

	return location, nil
}
