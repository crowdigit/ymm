package app

import (
	"github.com/crowdigit/ymm/db"
	"github.com/crowdigit/ymm/ydl"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=../mock/mock_app.go -package=mock github.com/crowdigit/ymm/app Application
type Application interface {
	DownloadPlaylist(url string) error
	DownloadSingle(url string) error
}

type ApplicationImpl struct {
	logger *zap.SugaredLogger
	ydl    ydl.YoutubeDL
	db     db.Database
}

func NewApplicationImpl(logger *zap.SugaredLogger, ydl ydl.YoutubeDL, db db.Database) Application {
	return ApplicationImpl{
		logger: logger,
		ydl:    ydl,
		db:     db,
	}
}

func (app ApplicationImpl) DownloadPlaylist(url string) error {
	// TODO retry failed downloads
	// TODO save download result
	// TODO configurable concurrent downloads
	// TODO run loudness scanner

	metadataBytes, err := app.ydl.PlaylistMetadata(url)
	if err != nil {
		return errors.Wrap(err, "failed to fetch playlist metadata")
	}

	metadata := make([]ydl.VideoMetadata, 0, len(metadataBytes))
	for _, metadatumBytes := range metadataBytes {
		metadatum := ydl.VideoMetadata{}
		if err := jsoniter.Unmarshal(metadatumBytes, &metadatum); err != nil {
			return errors.Wrap(err, "failed to unmarshal video metadata")
		}
		metadata = append(metadata, metadatum)
		if err := app.db.StoreMetadata(metadatum.ID, metadatumBytes); err != nil {
			return errors.Wrap(err, "failed to store video metadata")
		}
	}

	for _, metadatum := range metadata {
		if _, err := app.ydl.Download(metadatum); err != nil {
			return errors.Wrap(err, "failed to download video with metadata")
		}
	}

	return nil
}

func (app ApplicationImpl) DownloadSingle(url string) error {
	// TODO retry failed downloads
	// TODO save download result
	// TODO run loudness scanner

	metadataBytes, err := app.ydl.VideoMetadata(url)
	if err != nil {
		return errors.Wrap(err, "failed to fetch video metadata")
	}

	metadata := ydl.VideoMetadata{}
	if err := jsoniter.Unmarshal(metadataBytes, &metadata); err != nil {
		return errors.Wrap(err, "failed to unmarshal video metadata")
	}

	if err := app.db.StoreMetadata(metadata.ID, metadataBytes); err != nil {
		return errors.Wrap(err, "failed to store video metadata")
	}

	_, err = app.ydl.Download(metadata)
	if err != nil {
		return errors.Wrap(err, "failed to download video")
	}

	return nil
}
