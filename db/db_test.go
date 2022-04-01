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
}

func (s *DBTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *DBTestSuite) TestStoreMetadata() {
	id := "Ss-ba-g82-0"
	expectedMap := make(map[string]any)
	expectedMap["a"] = 123
	expectedMap["b"] = "asdf"
	expected, err := jsoniter.Marshal(expectedMap)
	s.Nil(err)

	db, err := db.NewDatabaseImpl(s.config, s.mockDb)
	s.Nil(err)
	s.Nil(db.StoreMetadata(id, expected))

	path := filepath.Join(s.config.MetadataDir, fmt.Sprintf("%s.json", id))
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	s.Nil(err)
	defer file.Close()

	got, err := io.ReadAll(file)
	s.Nil(err)

	s.Equal(expected, got)
}

func TestDBTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
