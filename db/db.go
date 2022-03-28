package db

//go:generate mockgen -destination=../mock/mock_database.go -package=mock github.com/crowdigit/ymm/db Database
type Database interface {
	StoreMetadata([]byte) error
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

func (db *DatabaseImpl) StoreMetadata(metadata []byte) error {
	panic("not implemented") // TODO: Implement
}
