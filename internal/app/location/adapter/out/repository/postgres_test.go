package repository_test

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.com/spacewalker/locations/internal/pkg/config"
	"gitlab.com/spacewalker/locations/internal/pkg/util"
	"path"
	"runtime"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type PostgresTestSuite struct {
	suite.Suite

	db *sql.DB
	m  *migrate.Migrate
}

func (s *PostgresTestSuite) SetupSuite() {
	var err error

	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(path.Dir(filename), "../../../../../..")

	cfg, err := config.LoadUserConfig(
		"user_test",
		path.Join(rootDir, "configs"),
	)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), cfg)

	dbSource := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	s.db, err = util.OpenDB(
		cfg.DBDriver,
		dbSource,
	)
	require.NoError(s.T(), err)

	migrationsPath := "file://" + path.Join(rootDir, "db/migrations/locations")

	driver, err := postgres.WithInstance(s.db, &postgres.Config{})
	require.NoError(s.T(), err)

	s.m, err = migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), s.m)
	require.NoError(s.T(), s.m.Up())
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
	require.NoError(s.T(), s.m.Down())
}

func TestPostgresTestSuite(t *testing.T) {
	// Skip tests when using "-short" flag.
	if testing.Short() {
		t.Skip("Skipping long-running tests")
	}
	suite.Run(t, new(PostgresTestSuite))
}
