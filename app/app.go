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
	"go.uber.org/zap"
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

func NewApplicationImpl(
	logger *zap.SugaredLogger,
	ydl ydl.YoutubeDL,
	loudness loudness.LoudnessScanner,
	jq jq.Jq,
	db db.Database,
	config ApplicationConfig,
) Application {
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
		return db.Uploader{}, fmt.Errorf("failed to insert uploader data: %w", err)
	}
	return uploader, nil
}

func getOrCreateUser(db_ db.Database, metadata ydl.VideoMetadata) (db.Uploader, error) {
	query := db.NewSelectUploaderQuery(db_.BunDB(), metadata.UploaderID)
	uploaders, err := db_.SelectUploader(query)
	if err != nil {
		return db.Uploader{}, fmt.Errorf("failed to select uploader data: %w", err)
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
		return fmt.Errorf("failed to fetch playlist metadata: %w", err)
	}

	metadataBytes1, err := app.jq.Slurp(metadataBytes0)
	if err != nil {
		return fmt.Errorf("failed to run jq command: %w", err)
	}

	metadataOriginal := []map[string]any{}
	if err := jsoniter.Unmarshal(metadataBytes1, &metadataOriginal); err != nil {
		return fmt.Errorf("failed to unmarshal metadata: %w", err)
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
		app.logger.Infow(
			"finisehd",
			"parsed",
			*parsed,
			"filtered",
			*filtered,
			"downloaded",
			*downloaded,
		)
	}(&parsed, &filtered, &downloaded)

	for _, metadatumOriginal := range metadataOriginal {
		parsed += 1
		metadatumBytes, err := jsoniter.Marshal(metadatumOriginal)
		if err != nil {
			return fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		metadatum := ydl.VideoMetadata{}
		if err := jsoniter.Unmarshal(metadatumBytes, &metadatum); err != nil {
			return fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		if _, exists := idSet[metadatum.ID]; exists {
			filtered += 1
			continue
		}

		query := db.NewSelectDownloadQuery(app.db.BunDB(), metadatum.ID)
		downloads, err := app.db.SelectDownload(query)
		if err != nil {
			return fmt.Errorf("failed to query download cache: %w", err)
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
				return fmt.Errorf("failed to make download directory: %w", err)
			}
			uploaderDirs[uploader.ID] = uploaderDir
		}

		if err := app.db.StoreMetadata(metadatum.ID, metadatumBytes); err != nil {
			return fmt.Errorf("failed to persist metadata into DB: %w", err)
		}

		metadataPairs = append(metadataPairs, MetadataPair{
			Metadata: metadatum,
			Bytes:    metadatumBytes,
		})
	}

	for _, pair := range metadataPairs {
		uploaderDirectory, exists := uploaderDirs[pair.Metadata.UploaderID]
		if !exists {
			panic("how?; this is definitely a bug")
		} else if _, err := app.ydl.Download(uploaderDirectory, pair.Metadata); err != nil {
			return fmt.Errorf("failed to download video: %w", err)
		}

		path := filepath.Join(uploaderDirectory, audioFilename(pair.Metadata.Filename))
		if err := app.loudness.Tag(path); err != nil {
			return fmt.Errorf("failed to generate Replaygain tag: %w", err)
		}

		query := db.NewInsertDownloadQuery(app.db.BunDB(), db.Download{ID: pair.Metadata.ID})
		if err := app.db.Insert(query); err != nil {
			return fmt.Errorf("failed to insert download result into DB: %w", err)
		}

		downloaded += 1
	}

	return nil
}

func (app ApplicationImpl) DownloadSingle(url string) error {
	metadataBytes, err := app.ydl.VideoMetadata(url)
	if err != nil {
		return fmt.Errorf("failed to fetch video metadata: %w", err)
	}

	metadata := ydl.VideoMetadata{}
	if err := jsoniter.Unmarshal(metadataBytes, &metadata); err != nil {
		return fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	selectDownload := db.NewSelectDownloadQuery(app.db.BunDB(), metadata.ID)
	downloads, err := app.db.SelectDownload(selectDownload)
	if err != nil {
		return fmt.Errorf("failed to query download cache: %w", err)
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
		return fmt.Errorf("failed to persist metadata into DB: %w", err)
	}

	downloadDir := filepath.Join(app.config.DownloadRootDir, uploader.Directory)
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return fmt.Errorf("failed to make download directory: %w", err)
	}

	_, err = app.ydl.Download(downloadDir, metadata)
	if err != nil {
		return fmt.Errorf("failed to download video: %w", err)
	}

	path := filepath.Join(
		app.config.DownloadRootDir,
		uploader.Directory,
		audioFilename(metadata.Filename),
	)
	if err := app.loudness.Tag(path); err != nil {
		return fmt.Errorf("failed to generate Replaygain tag: %w", err)
	}

	insertDownload := db.NewInsertDownloadQuery(app.db.BunDB(), db.Download{ID: metadata.ID})
	if err := app.db.Insert(insertDownload); err != nil {
		return fmt.Errorf("failed to insert download result into DB: %w", err)
	}

	return nil
}
