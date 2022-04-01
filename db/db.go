package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type Uploader struct {
	bun.BaseModel `bun:"table:uploaders"`

	ID        string `bun:",pk"`
	URL       string
	Name      string
	Directory string
}

//go:generate mockgen -destination=../mock/mock_database.go -package=mock github.com/crowdigit/ymm/db Database
type Database interface {
	StoreMetadata(id string, metadata []byte) error
	Close()
}

type DatabaseConfig struct {
	DatabaseFile string
	MetadataDir  string
}

type DatabaseImpl struct {
	config DatabaseConfig
	sqldb  *sql.DB
}

func NewDatabaseImpl(config DatabaseConfig) (*DatabaseImpl, error) {
	sqldb, err := sql.Open(sqliteshim.ShimName, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open sqlite DB")
	}

	return &DatabaseImpl{
		config: config,
		sqldb:  sqldb,
	}, nil
}

func (db *DatabaseImpl) StoreMetadata(id string, metadata []byte) error {
	path := filepath.Join(db.config.MetadataDir, fmt.Sprintf("%s.json", id))

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to open metadata file")
	}
	defer file.Close()

	if _, err := file.Write(metadata); err != nil {
		return errors.Wrap(err, "failed to write metadata file")
	}

	return nil
}

func (db *DatabaseImpl) Close() {
	db.sqldb.Close()
}
