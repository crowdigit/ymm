package db

import (
	"context"
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
	URL       string `bun:",notnull"`
	Name      string `bun:",notnull"`
	Directory string `bun:",notnull"`
}

type Download struct {
	bun.BaseModel `bun:"table:downloads"`

	ID string `bun:",pk"`
}

//go:generate mockgen -destination=../mock/mock_database.go -package=mock github.com/crowdigit/ymm/db Database
type Database interface {
	StoreMetadata(id string, metadata []byte) error
	InsertUploader(*bun.InsertQuery) error
	SelectUploader(*bun.SelectQuery) ([]Uploader, error)
	InsertDownload(*bun.InsertQuery) error
	SelectDownload(*bun.SelectQuery) ([]Download, error)
	BunDB() *bun.DB
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

	ctx := context.Background()
	if _, err := bundb.NewCreateTable().IfNotExists().Table("uploaders").Model(&Uploader{}).Exec(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to create uploaders table")
	}

	if _, err := bundb.NewCreateTable().IfNotExists().Table("downloads").Model(&Download{}).Exec(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to create downloads table")
	}

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

func (db *DatabaseImpl) InsertUploader(query *bun.InsertQuery) error {
	ctx := context.Background()
	_, err := query.Exec(ctx)
	return errors.Wrap(err, "failed to execute insert query")
}

func (db *DatabaseImpl) SelectUploader(query *bun.SelectQuery) ([]Uploader, error) {
	var uploaders []Uploader
	ctx := context.Background()
	err := query.Model(&uploaders).Scan(ctx)
	return uploaders, errors.Wrap(err, "failed to query scan query result")
}

func (db *DatabaseImpl) InsertDownload(query *bun.InsertQuery) error {
	ctx := context.Background()
	_, err := query.Exec(ctx)
	return errors.Wrap(err, "failed to execute insert query")
}

func (db *DatabaseImpl) SelectDownload(query *bun.SelectQuery) ([]Download, error) {
	var downloads []Download
	ctx := context.Background()
	err := query.Model(&downloads).Scan(ctx)
	return downloads, errors.Wrap(err, "failed to query download result")
}

func (db *DatabaseImpl) BunDB() *bun.DB {
	return db.bundb
}
