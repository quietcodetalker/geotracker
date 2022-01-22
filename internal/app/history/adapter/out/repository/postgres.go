package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"gitlab.com/spacewalker/locations/internal/app/history/core/domain"
	"gitlab.com/spacewalker/locations/internal/app/history/core/port"
	"gitlab.com/spacewalker/locations/internal/pkg/geo"
)

const (
	// RecordsTable contains name of the records database table.
	RecordsTable = "records"

	constraintRecordsALongitudeValid = "records_a_longitude_valid"
	constraintRecordsALatitudeValid  = "records_a_latitude_valid"
	constraintRecordsBLongitudeValid = "records_b_longitude_valid"
	constraintRecordsBLatitudeValid  = "records_b_latitude_valid"
)

type postgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates postgres repository instance and return its pointer.
func NewPostgresRepository(db *sql.DB) port.HistoryRepository {
	return &postgresRepository{db: db}
}

var addRecordQuery = fmt.Sprintf(
	`
INSERT INTO %s
(user_id, a, b)
VALUES ($1, $2, $3)
RETURNING id, user_id, a, b, created_at, updated_at
`,
	RecordsTable,
)

// AddRecord adds a record with to records table and returns it.
func (r postgresRepository) AddRecord(ctx context.Context, req port.HistoryRepositoryAddRecordRequest) (domain.Record, error) {
	var record domain.Record
	var a, b geo.PostgresPoint

	if err := r.db.QueryRowContext(ctx, addRecordQuery, req.UserID, geo.PostgresPoint(req.A), geo.PostgresPoint(req.B)).Scan(
		&record.ID,
		&record.UserID,
		&a,
		&b,
		&record.CreatedAt,
		&record.UpdatedAt,
	); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Constraint {
			case constraintRecordsALongitudeValid:
				return domain.Record{}, &port.InvalidLocationError{
					Violations: []port.InvalidLocationErrorViolation{
						{
							Subject: "a.longitude",
							Value:   req.A.Longitude(),
						},
					},
				}
			case constraintRecordsALatitudeValid:
				return domain.Record{}, &port.InvalidLocationError{
					Violations: []port.InvalidLocationErrorViolation{
						{
							Subject: "a.latitude",
							Value:   req.A.Latitude(),
						},
					},
				}
			case constraintRecordsBLongitudeValid:
				return domain.Record{}, &port.InvalidLocationError{
					Violations: []port.InvalidLocationErrorViolation{
						{
							Subject: "b.longitude",
							Value:   req.B.Longitude(),
						},
					},
				}
			case constraintRecordsBLatitudeValid:
				return domain.Record{}, &port.InvalidLocationError{
					Violations: []port.InvalidLocationErrorViolation{
						{
							Subject: "b.latitude",
							Value:   req.B.Latitude(),
						},
					},
				}
			}
		}
		return domain.Record{}, err
	}

	record.A = geo.Point(a)
	record.B = geo.Point(b)

	return record, nil
}

var getDistanceQuery = fmt.Sprintf(
	`
SELECT coalesce(SUM(a <@> b), 0.00) * 1609.344
FROM %s
WHERE user_id = $1 AND created_at >= $2 AND created_at <= $3
`,
	RecordsTable,
)

// GetDistance finds distance that users got through in a period of time.
func (r postgresRepository) GetDistance(ctx context.Context, req port.HistoryRepositoryGetDistanceRequest) (float64, error) {
	var distance float64
	if err := r.db.QueryRowContext(ctx, getDistanceQuery, req.UserID, req.From, req.To).Scan(&distance); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return distance, nil
}
