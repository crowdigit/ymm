package db

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

//go:generate mockgen -destination=../mock/mock_database.go -package=mock github.com/crowdigit/ymm/db Database
type Database interface {
	StoreMetadata(id string, metadata []byte) error
}

type DatabaseConfig struct {
	DatabaseFile string
	MetadataDir  string
}

type DatabaseImpl struct {
	config DatabaseConfig
}

func NewDatabaseImpl(config DatabaseConfig) *DatabaseImpl {
	return &DatabaseImpl{
		config: config,
	}
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
