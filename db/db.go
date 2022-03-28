package db

//go:generate mockgen -destination=../mock/mock_database.go -package=mock github.com/crowdigit/ymm/db Database
type Database interface {
	StoreMetadata([]byte) error
}
