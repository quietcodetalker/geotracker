package repository_test

import (
	"fmt"

	"github.com/stretchr/testify/require"
	"gitlab.com/spacewalker/geotracker/internal/app/history/adapter/out/repository"
	"gitlab.com/spacewalker/geotracker/internal/app/history/core/domain"
	"gitlab.com/spacewalker/geotracker/internal/pkg/geo"
)

var seedRecordsQuery = fmt.Sprintf(
	`
INSERT INTO %s
(user_id, a, b, timestamp)
VALUES ($1, $2, $3, $4)
RETURNING user_id, a, b, timestamp
`,
	repository.RecordsTable,
)

func (s *PostgresTestSuite) seedRecords(args []domain.Record) []domain.Record {
	records := make([]domain.Record, 0, len(args))

	tx, err := s.db.Begin()
	require.NoError(s.T(), err)

	defer tx.Rollback()

	stmt, err := tx.Prepare(seedRecordsQuery)
	require.NoError(s.T(), err)
	defer stmt.Close()

	for _, arg := range args {
		var record domain.Record
		var a, b geo.PostgresPoint

		err := stmt.QueryRow(arg.UserID, geo.PostgresPoint(arg.A), geo.PostgresPoint(arg.B), arg.Timestamp).Scan(
			&record.ID,
			&a,
			&b,
			&record.Timestamp,
		)
		require.NoError(s.T(), err)

		record.A = geo.Point(a)
		record.B = geo.Point(b)
		records = append(records, record)
	}

	err = tx.Commit()
	require.NoError(s.T(), err)

	return records
}
