package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"gitlab.com/spacewalker/geotracker/internal/pkg/util/testutil"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.com/spacewalker/geotracker/internal/pkg/util"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type PostgresTestSuite struct {
	suite.Suite

	db        *sql.DB
	m         *migrate.Migrate
	container *testutil.Container
}

func (s *PostgresTestSuite) SetupSuite() {
	var err error

	dbUser := testutil.RandomString(10, 10, testutil.CharacterSetAlphabet)
	dbPassword := testutil.RandomString(10, 10, testutil.CharacterSetAlphabet)
	dbName := testutil.RandomString(10, 10, testutil.CharacterSetAlphabet)

	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(path.Dir(filename), "../../../../../..")

	// Setup postgres in a docker container.
	cancelCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	s.container, err = testutil.SetupPostgres(cancelCtx, testutil.PostgresConfig{
		User:     dbUser,
		Password: dbPassword,
		DBName:   dbName,
	})
	if err != nil {
		s.T().Fatal(err)
	}

	// Connect to database.
	s.db, err = util.OpenDB(
		"postgres",
		s.container.URI,
	)
	if err != nil {
		s.T().Fatal(err)
	}

	if err := s.db.Ping(); err != nil {
		s.T().Fatal(err)
	}

	// Run migrations.
	migrationsPath := "file://" + path.Join(rootDir, "db/migrations/locations")

	driver, err := postgres.WithInstance(s.db, &postgres.Config{})
	if err != nil {
		s.T().Fatal(err)
	}

	s.m, err = migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		s.T().Fatal(err)
	}

	err = s.m.Up()
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *PostgresTestSuite) TearDownTest() {
	query := `
	SELECT tablename
	FROM pg_catalog.pg_tables
	WHERE schemaname != 'pg_catalog' AND
				schemaname != 'information_schema' AND
				tablename != 'schema_migrations'
	`
	rows, err := s.db.Query(query)
	require.NoError(s.T(), err)
	defer rows.Close()

	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		require.NoError(s.T(), err)

		truncateQuery := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tableName)
		_, err = s.db.Exec(truncateQuery)
		require.NoError(s.T(), err)
	}

	err = rows.Close()
	require.NoError(s.T(), err)
}

func (s *PostgresTestSuite) TearDownSuite() {
	var err error

	err = s.m.Down()
	if err != nil {
		// TODO: Log err
	}
	err = s.db.Close()
	if err != nil {
		// TODO: Log err
	}
	cancelCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err = s.container.Terminate(cancelCtx)
	if err != nil {
		s.T().Fatal(err)
	}
}

func TestPostgresTestSuite(t *testing.T) {
	// Skip tests when using "-short" flag.
	if testing.Short() {
		t.Skip("Skipping long-running tests")
	}
	suite.Run(t, new(PostgresTestSuite))
}
