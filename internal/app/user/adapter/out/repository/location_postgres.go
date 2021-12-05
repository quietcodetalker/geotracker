package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"gitlab.com/spacewalker/locations/internal/app/user/core/domain"
	"gitlab.com/spacewalker/locations/internal/app/user/core/port"
)

var setLocationQuery = fmt.Sprintf(
	`
INSERT INTO %s
(user_id, latitude, longitude)
VALUES ($1, $2, $3)
ON CONFLICT ON CONSTRAINT locations_pkey 
DO
	UPDATE SET latitude = EXCLUDED.latitude, longitude = EXCLUDED.longitude
RETURNING user_id, latitude, longitude, created_at, updated_at
`,
	LocationTable,
)

// SetLocation sets user's location by given user ID.
// Returns updated location entity.
func (q *postgresQueries) SetLocation(ctx context.Context, arg port.SetLocationArg) (domain.Location, error) {
	var location domain.Location
	if err := q.db.QueryRowContext(ctx, setLocationQuery, arg.UserID, arg.Latitude, arg.Longitude).Scan(
		&location.UserID,
		&location.Latitude,
		&location.Longitude,
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
							Value:   arg.Latitude,
						},
					},
				}
			case ConstraintLocationsLongitudeValid:
				return domain.Location{}, &port.InvalidLocationError{
					Violations: []port.InvalidLocationErrorViolation{
						{
							Subject: "longitude",
							Value:   arg.Longitude,
						},
					},
				}
			}
		}
		return domain.Location{}, err
	}

	return location, nil
}
