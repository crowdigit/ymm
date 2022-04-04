package db_test

import (
	"database/sql"
	"testing"

	"github.com/crowdigit/ymm/db"
	"github.com/stretchr/testify/suite"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type QueryTestSuite struct {
	suite.Suite

	sqldb  *sql.DB
	bundb  *bun.DB
	dbImpl *db.DatabaseImpl
}

func (s *QueryTestSuite) SetupTest() {
	sqldb, err := sql.Open(sqliteshim.ShimName, ":memory:")
	s.Nil(err)
	s.sqldb = sqldb

	dbImpl, err := db.NewDatabaseImpl(db.DatabaseConfig{}, sqldb)
	s.Nil(err)
	s.dbImpl = dbImpl
	s.bundb = dbImpl.BunDB()

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
	s.Nil(s.bundb.Close())
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

	query := db.NewGetUploaderQuery(s.bundb, uploader.ID)
	uploaders, err := s.dbImpl.SelectUploader(query)
	s.Nil(err)
	s.Len(uploaders, 1)
	s.Equal(uploader, uploaders[0])
}

func TestQueryTestSuite(t *testing.T) {
	suite.Run(t, new(QueryTestSuite))
}
