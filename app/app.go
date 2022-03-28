package app

import (
	"encoding/json"

	"github.com/crowdigit/ymm/ydl"
	"github.com/pkg/errors"
)

//go:generate mockgen -destination=../mock/mock_app.go -package=mock github.com/crowdigit/ymm/app Application
type Application interface {
	DownloadPlaylist(url string) error
	DownloadSingle(url string) error
}

type ApplicationImpl struct {
	ydl ydl.YoutubeDL
}

func NewApplicationImpl(ydl ydl.YoutubeDL) Application {
	return ApplicationImpl{
		ydl: ydl,
	}
}

func (app ApplicationImpl) DownloadPlaylist(url string) error {
	metadata, err := app.ydl.PlaylistMetadata(url)
	if err != nil {
		return errors.Wrap(err, "failed to fetch playlist metadata")
	}

	// TODO configurable concurrent downloads
	for _, metadatumBytes := range metadata {
		// TODO retry failed downloads
		metadatum := ydl.VideoMetadata{}
		if err := json.Unmarshal(metadatumBytes, &metadatum); err != nil {
			return err
		}
		if _, err := app.ydl.Download(metadatum); err != nil {
			return err
		}
	}

	return nil
}

func (app ApplicationImpl) DownloadSingle(url string) error {
	metadataBytes, err := app.ydl.VideoMetadata(url)
	if err != nil {
		return errors.Wrap(err, "failed to fetch video metadata")
	}

	metadata := ydl.VideoMetadata{}
	if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
		return err
	}

	_, err = app.ydl.Download(metadata)
	if err != nil {
		return errors.Wrap(err, "failed to download video")
	}

	return nil
}
