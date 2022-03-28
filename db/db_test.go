package db_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/crowdigit/ymm/db"
	"github.com/golang/mock/gomock"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/suite"
)

type DBTestSuite struct {
	suite.Suite

	mockCtrl *gomock.Controller
	dataDir  string
	config   db.DatabaseConfig
}

func (s *DBTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.dataDir = s.T().TempDir()
	s.config = db.DatabaseConfig{
		DatabaseFile: filepath.Join(s.dataDir, "db.sql"),
		MetadataDir:  filepath.Join(s.dataDir, "metadata"),
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

	db := db.NewDatabaseImpl(s.config)
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
