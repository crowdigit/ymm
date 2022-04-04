package db_test

import (
	"database/sql"
	"testing"

	"github.com/crowdigit/ymm/db"
	"github.com/stretchr/testify/suite"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type QueryTestSuite struct {
	suite.Suite

	sqldb *sql.DB
	db    *db.DatabaseImpl
}

func (s *QueryTestSuite) SetupTest() {
	sqldb, err := sql.Open(sqliteshim.ShimName, ":memory:")
	s.Nil(err)
	s.sqldb = sqldb

	dbImpl, err := db.NewDatabaseImpl(db.DatabaseConfig{}, sqldb)
	s.Nil(err)
	s.db = dbImpl

	_, err = s.sqldb.Exec(
		`CREATE TABLE IF NOT EXISTS uploaders
		( id        TEXT PRIMARY KEY NOT NULL
		, url       TEXT             NOT NULL
		, name      TEXT             NOT NULL
		, directory TEXT             NOT NULL
		)`)
	s.Nil(err)
}

func (s *QueryTestSuite) TearDownTest() {
	s.Nil(s.sqldb.Close())
}

func (s *QueryTestSuite) TestSelectUser() {
	uploader := db.Uploader{
		ID:        "uploader id",
		URL:       "uploader url",
		Name:      "uploader name",
		Directory: "uploader directory",
	}
	_, err := s.sqldb.Exec(
		`INSERT INTO uploaders
		( id, url, name, directory )
		VALUES
		( ?, ?, ?, ? )`,
		uploader.ID,
		uploader.URL,
		uploader.Name,
		uploader.Directory)
	s.Nil(err)

	query := db.NewSelectUploaderQuery(s.db.BunDB(), uploader.ID)
	uploaders, err := s.db.SelectUploader(query)
	s.Nil(err)
	s.Len(uploaders, 1)
	s.Equal(uploader, uploaders[0])
}

func (s *QueryTestSuite) TestInsertUser() {
	expected := db.Uploader{
		ID:        "uploader id",
		URL:       "uploader url",
		Name:      "uploader name",
		Directory: "uploader directory",
	}

	query := db.NewInsertUploaderQuery(s.db.BunDB(), expected)
	s.Nil(s.db.InsertUploader(query))

	rows, err := s.sqldb.Query(
		`SELECT * FROM uploaders WHERE id = ?`,
		expected.ID,
	)
	s.Nil(err)
	s.True(rows.Next())

	got := db.Uploader{}
	s.Nil(
		rows.Scan(
			&got.ID,
			&got.URL,
			&got.Name,
			&got.Directory))
	s.Equal(expected, got)
}

func TestQueryTestSuite(t *testing.T) {
	suite.Run(t, new(QueryTestSuite))
}
