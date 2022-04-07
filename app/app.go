package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/crowdigit/ymm/db"
	"github.com/crowdigit/ymm/jq"
	"github.com/crowdigit/ymm/loudness"
	"github.com/crowdigit/ymm/ydl"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	ERR_MARSHAL_METADATA       = "failed to marshal video metadata"
	ERR_UNMARSHAL_METADATA     = "failed to unmarshal video metadata"
	ERR_STORE_VIDEO_METADATA   = "failed to store video metadata"
	ERR_DOWNLOAD_VIDEO         = "failed to download video"
	ERR_TAG_LOUDNESS           = "failed to tag loudness"
	ERR_MAKE_DOWNLOAD_DIR      = "failed to make download directory"
	ERR_INSERT_DOWNLOAD_RESULT = "failed to insert download result"
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
	jq       jq.Jq
	db       db.Database
	config   ApplicationConfig
}

func NewApplicationImpl(logger *zap.SugaredLogger, ydl ydl.YoutubeDL, loudness loudness.LoudnessScanner, jq jq.Jq, db db.Database, config ApplicationConfig) Application {
	return ApplicationImpl{
		logger:   logger,
		ydl:      ydl,
		loudness: loudness,
		jq:       jq,
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
	if err := db_.Insert(query); err != nil {
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

type MetadataPair struct {
	Metadata ydl.VideoMetadata
	Bytes    []byte
}

func (app ApplicationImpl) DownloadPlaylist(url string) error {
	// TODO retry failed downloads
	// TODO save download result
	// TODO configurable concurrent downloads

	metadataBytes0, err := app.ydl.VideoMetadata(url)
	if err != nil {
		return errors.Wrap(err, "failed to fetch playlist metadata")
	}

	metadataBytes1, err := app.jq.Slurp(metadataBytes0)
	if err != nil {
		return errors.Wrap(err, "failed to run jq command")
	}

	metadataOriginal := []map[string]any{}
	if err := jsoniter.Unmarshal(metadataBytes1, &metadataOriginal); err != nil {
		return errors.Wrap(err, ERR_UNMARSHAL_METADATA)
	}

	uploaderDirs := make(map[string]string)

	metadataPairs := make(
		[]MetadataPair,
		0,
		len(metadataOriginal))

	parsed := 0
	filtered := 0
	downloaded := 0

	// DB downloads table has unitque constraint, so it needs to filter the same
	// contents from a playlist
	idSet := make(map[string]struct{})

	defer func(parsed *int, filtered *int, downloaded *int) {
		app.logger.Infow("finisehd", "parsed", *parsed, "filtered", *filtered, "downloaded", *downloaded)
	}(&parsed, &filtered, &downloaded)

	for _, metadatumOriginal := range metadataOriginal {
		parsed += 1
		metadatumBytes, err := jsoniter.Marshal(metadatumOriginal)
		if err != nil {
			return errors.Wrap(err, ERR_MARSHAL_METADATA)
		}

		metadatum := ydl.VideoMetadata{}
		if err := jsoniter.Unmarshal(metadatumBytes, &metadatum); err != nil {
			return errors.Wrap(err, ERR_UNMARSHAL_METADATA)
		}

		if _, exists := idSet[metadatum.ID]; exists {
			filtered += 1
			continue
		}

		query := db.NewSelectDownloadQuery(app.db.BunDB(), metadatum.ID)
		downloads, err := app.db.SelectDownload(query)
		if err != nil {
			errors.Wrap(err, "failed to query download cache")
		}
		if len(downloads) > 0 {
			filtered += 1
			continue
		}

		idSet[metadatum.ID] = struct{}{}

		uploader, err := getOrCreateUser(app.db, metadatum)
		if err != nil {
			return err
		}

		if _, exists := uploaderDirs[uploader.ID]; !exists {
			uploaderDir := filepath.Join(app.config.DownloadRootDir, uploader.Directory)
			if err := os.MkdirAll(uploaderDir, 0755); err != nil {
				return errors.Wrap(err, ERR_MAKE_DOWNLOAD_DIR)
			}
			uploaderDirs[uploader.ID] = uploaderDir
		}

		if err := app.db.StoreMetadata(metadatum.ID, metadatumBytes); err != nil {
			return errors.Wrap(err, ERR_STORE_VIDEO_METADATA)
		}

		metadataPairs = append(metadataPairs, MetadataPair{
			Metadata: metadatum,
			Bytes:    metadatumBytes,
		})
	}

	for _, pair := range metadataPairs {
		uploaderDirectory, exists := uploaderDirs[pair.Metadata.UploaderID]
		if !exists {
			return errors.New("how?; this is definitely a bug")
		} else if _, err := app.ydl.Download(uploaderDirectory, pair.Metadata); err != nil {
			return errors.Wrap(err, ERR_DOWNLOAD_VIDEO)
		}

		path := filepath.Join(uploaderDirectory, audioFilename(pair.Metadata.Filename))
		if err := app.loudness.Tag(path); err != nil {
			return errors.Wrap(err, ERR_TAG_LOUDNESS)
		}

		query := db.NewInsertDownloadQuery(app.db.BunDB(), db.Download{ID: pair.Metadata.ID})
		if err := app.db.Insert(query); err != nil {
			return errors.Wrap(err, ERR_INSERT_DOWNLOAD_RESULT)
		}

		downloaded += 1
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

	selectDownload := db.NewSelectDownloadQuery(app.db.BunDB(), metadata.ID)
	downloads, err := app.db.SelectDownload(selectDownload)
	if err != nil {
		return errors.Wrap(err, "failed to query download cache")
	}
	if len(downloads) > 0 {
		app.logger.Infof("%s is already downloaded", metadata.ID)
		return nil
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

	insertDownload := db.NewInsertDownloadQuery(app.db.BunDB(), db.Download{ID: metadata.ID})
	if err := app.db.Insert(insertDownload); err != nil {
		return errors.Wrap(err, ERR_INSERT_DOWNLOAD_RESULT)
	}

	return nil
}
