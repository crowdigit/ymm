package db_test

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/crowdigit/ymm/db"
	"github.com/golang/mock/gomock"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/suite"
)

type DBTestSuite struct {
	suite.Suite

	mockCtrl *gomock.Controller
	mockSql  sqlmock.Sqlmock
	mockDb   *sql.DB
	dataDir  string
	config   db.DatabaseConfig

	db *db.DatabaseImpl
}

func (s *DBTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	mockDb, mockSql, err := sqlmock.New()
	s.Nil(err)
	s.mockDb, s.mockSql = mockDb, mockSql
	s.dataDir = s.T().TempDir()
	s.config = db.DatabaseConfig{
		MetadataDir: filepath.Join(s.dataDir, "metadata"),
	}
	s.Nil(os.Mkdir(s.config.MetadataDir, 0755))

	s.mockSql.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(0, 0))
	s.mockSql.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(0, 0))
	db, err := db.NewDatabaseImpl(s.config, s.mockDb)
	s.Nil(err)
	s.db = db
}

func (s *DBTestSuite) TearDownTest() {
	if err := s.mockSql.ExpectationsWereMet(); err != nil {
		s.Error(err)
	}
	s.mockCtrl.Finish()
}

func (s *DBTestSuite) TestStoreMetadata() {
	id := "Ss-ba-g82-0"
	expectedMap := make(map[string]any)
	expectedMap["a"] = 123
	expectedMap["b"] = "asdf"
	expected, err := jsoniter.Marshal(expectedMap)
	s.Nil(err)

	s.Nil(s.db.StoreMetadata(id, expected))

	path := filepath.Join(s.config.MetadataDir, fmt.Sprintf("%s.json", id))
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	s.Nil(err)
	defer file.Close()

	got, err := io.ReadAll(file)
	s.Nil(err)

	s.Equal(expected, got)
}

func (s *DBTestSuite) TestInsertUser() {
	uploader := db.Uploader{
		ID:        "uploader id",
		URL:       "uploader url",
		Name:      "uploader name",
		Directory: "uploader directory",
	}
	s.mockSql.
		ExpectExec("INSERT").
		WillReturnResult(sqlmock.NewResult(1, 1))
	query := db.NewInsertUploaderQuery(s.db.BunDB(), uploader)
	s.Nil(s.db.Insert(query))
}

func (s *DBTestSuite) TestSelectSingleUser() {
	uploader := db.Uploader{
		ID:        "uploader id",
		URL:       "uploader url",
		Name:      "uploader name",
		Directory: "uploader directory",
	}
	rows := sqlmock.NewRows([]string{"id", "url", "name", "directory"}).
		AddRow(uploader.ID, uploader.URL, uploader.Name, uploader.Directory)
	s.mockSql.
		ExpectQuery("SELECT").
		WillReturnRows(rows)
	query := db.NewSelectUploaderQuery(s.db.BunDB(), uploader.ID)
	uploaders, err := s.db.SelectUploader(query)
	s.Nil(err)
	s.Len(uploaders, 1)
	s.Equal(uploader, uploaders[0])
}

func TestDBTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
