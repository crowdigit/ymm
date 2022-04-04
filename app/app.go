package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/crowdigit/ymm/db"
	"github.com/crowdigit/ymm/loudness"
	"github.com/crowdigit/ymm/ydl"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	ERR_UNMARSHAL_METADATA   = "failed to unmarshal video metadata"
	ERR_STORE_VIDEO_METADATA = "failed to store video metadata"
	ERR_DOWNLOAD_VIDEO       = "failed to download video"
	ERR_TAG_LOUDNESS         = "failed to tag loudness"
	ERR_MAKE_DOWNLOAD_DIR    = "failed to make download directory"
)

//go:generate mockgen -destination=../mock/mock_app.go -package=mock github.com/crowdigit/ymm/app Application
type Application interface {
	DownloadPlaylist(url string) error
	DownloadSingle(url string) error
}

type ApplicationConfig struct {
	DownloadRootDir string
}

type ApplicationImpl struct {
	logger   *zap.SugaredLogger
	ydl      ydl.YoutubeDL
	loudness loudness.LoudnessScanner
	db       db.Database
	config   ApplicationConfig
}

func NewApplicationImpl(logger *zap.SugaredLogger, ydl ydl.YoutubeDL, loudness loudness.LoudnessScanner, db db.Database, config ApplicationConfig) Application {
	return ApplicationImpl{
		logger:   logger,
		ydl:      ydl,
		loudness: loudness,
		db:       db,
		config:   config,
	}
}

func audioFilename(videoFilename string) string {
	return fmt.Sprintf(
		"%s.mp3",
		strings.TrimSuffix(videoFilename, filepath.Ext(videoFilename)))
}

func insertUser(db_ db.Database, metadata ydl.VideoMetadata) (db.Uploader, error) {
	uploader := db.Uploader{
		ID:        metadata.UploaderID,
		URL:       metadata.UploaderURL,
		Name:      metadata.Uploader,
		Directory: metadata.UploaderID,
	}
	query := db.NewInsertUploaderQuery(db_.BunDB(), uploader)
	if err := db_.InsertUploader(query); err != nil {
		return db.Uploader{}, errors.Wrap(err, "failed to insert uploader data")
	}
	return uploader, nil
}

func getOrCreateUser(db_ db.Database, metadata ydl.VideoMetadata) (db.Uploader, error) {
	query := db.NewSelectUploaderQuery(db_.BunDB(), metadata.UploaderID)
	uploaders, err := db_.SelectUploader(query)
	if err != nil {
		return db.Uploader{}, errors.Wrap(err, "failed to select uploader data")
	}

	if len(uploaders) > 0 {
		return uploaders[0], nil
	}

	return insertUser(db_, metadata)
}

func (app ApplicationImpl) DownloadPlaylist(url string) error {
	// TODO retry failed downloads
	// TODO save download result
	// TODO configurable concurrent downloads

	metadataBytes, err := app.ydl.PlaylistMetadata(url)
	if err != nil {
		return errors.Wrap(err, "failed to fetch playlist metadata")
	}

	metadata := make([]ydl.VideoMetadata, 0, len(metadataBytes))
	uploaderDirs := make(map[string]string)
	for _, metadatumBytes := range metadataBytes {
		metadatum := ydl.VideoMetadata{}
		if err := jsoniter.Unmarshal(metadatumBytes, &metadatum); err != nil {
			return errors.Wrap(err, "failed to unmarshal video metadata")
		}

		uploader, err := getOrCreateUser(app.db, metadatum)
		if err != nil {
			return err
		}

		downloadDir := filepath.Join(app.config.DownloadRootDir, uploader.Directory)
		if err := os.MkdirAll(downloadDir, 0755); err != nil {
			return errors.Wrap(err, ERR_MAKE_DOWNLOAD_DIR)
		}
		uploaderDirs[uploader.ID] = downloadDir

		metadata = append(metadata, metadatum)
		if err := app.db.StoreMetadata(metadatum.ID, metadatumBytes); err != nil {
			return errors.Wrap(err, ERR_UNMARSHAL_METADATA)
		}
	}

	for _, metadatum := range metadata {
		uploaderDirectory := uploaderDirs[metadatum.UploaderID]
		if _, err := app.ydl.Download(uploaderDirectory, metadatum); err != nil {
			return errors.Wrap(err, ERR_DOWNLOAD_VIDEO)
		}

		path := filepath.Join(uploaderDirectory, audioFilename(metadatum.Filename))
		if err := app.loudness.Tag(path); err != nil {
			return errors.Wrap(err, ERR_TAG_LOUDNESS)
		}
	}

	return nil
}

func (app ApplicationImpl) DownloadSingle(url string) error {
	// TODO retry failed downloads
	// TODO save download result

	metadataBytes, err := app.ydl.VideoMetadata(url)
	if err != nil {
		return errors.Wrap(err, "failed to fetch video metadata")
	}

	metadata := ydl.VideoMetadata{}
	if err := jsoniter.Unmarshal(metadataBytes, &metadata); err != nil {
		return errors.Wrap(err, ERR_UNMARSHAL_METADATA)
	}

	uploader, err := getOrCreateUser(app.db, metadata)
	if err != nil {
		return err
	}

	if err := app.db.StoreMetadata(metadata.ID, metadataBytes); err != nil {
		return errors.Wrap(err, ERR_STORE_VIDEO_METADATA)
	}

	downloadDir := filepath.Join(app.config.DownloadRootDir, uploader.Directory)
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return errors.Wrap(err, ERR_MAKE_DOWNLOAD_DIR)
	}

	_, err = app.ydl.Download(downloadDir, metadata)
	if err != nil {
		return errors.Wrap(err, ERR_DOWNLOAD_VIDEO)
	}

	path := filepath.Join(app.config.DownloadRootDir, uploader.Directory, audioFilename(metadata.Filename))
	if err := app.loudness.Tag(path); err != nil {
		return errors.Wrap(err, ERR_TAG_LOUDNESS)
	}

	return nil
}
