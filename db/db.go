package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
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
	SetUploader(bun.Query) error
	GetUploader(bun.Query) ([]Uploader, error)
}

type DatabaseConfig struct {
	MetadataDir string
}

type DatabaseImpl struct {
	config DatabaseConfig
	bundb  *bun.DB
}

func NewDatabaseImpl(config DatabaseConfig, sqldb *sql.DB) (*DatabaseImpl, error) {
	bundb := bun.NewDB(sqldb, sqlitedialect.New())

	return &DatabaseImpl{
		config: config,
		bundb:  bundb,
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

func (db *DatabaseImpl) SetUploader(query bun.Query) error {
	panic("not implemented") // TODO: Implement
}

func (db *DatabaseImpl) GetUploader(query bun.Query) ([]Uploader, error) {
	panic("not implemented") // TODO: Implement
}

func (db *DatabaseImpl) BunDB() *bun.DB {
	return db.bundb
}
