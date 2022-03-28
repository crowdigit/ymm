package db

import "github.com/crowdigit/ymm/ydl"

//go:generate mockgen -destination=../mock/mock_database.go -package=mock github.com/crowdigit/ymm/db Database
type Database interface {
	StoreMetadata(ydl.VideoMetadata) error
}
