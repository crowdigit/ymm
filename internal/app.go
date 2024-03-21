package internal

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/crowdigit/exec"
)

func downloadPath(metadata JqMetadata) (string, error) {
	filename := strings.TrimSuffix(
		metadata.Filename,
		filepath.Ext(metadata.Filename),
	) + ".mp3"
	path, err := xdg.DataFile(filepath.Join("ymm", metadata.UploaderID, filename))
	if err != nil {
		return "", fmt.Errorf("failed to determine download location: %w", err)
	}
	return path, nil
}

func DownloadSingle(cp exec.CommandProvider, config Config, url string) error {
	jqMetadata, err := fetchMetadata(
		cp,
		config.Command.Metadata.Youtube,
		config.Command.Metadata.JSON,
		url,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch media metadata: %w", err)
	}
	if len(jqMetadata) != 1 {
		// TODO error with context
		return fmt.Errorf("fetched zero metadata")
	}

	format, err := selectFormat(jqMetadata[0])
	if err != nil {
		return fmt.Errorf("failed to validate available media formats: %w", err)
	}
	config.Command.Download.Youtube.ReplacePlaceholder("<format>", format)

	// TODO persist metadata into DB

	path, err := downloadPath(jqMetadata[0])
	if err != nil {
		return fmt.Errorf("failed to build download path: %w", err)
	}

	if err := downloadVideo(cp, config.Command.Download.Youtube, url, format, filepath.Dir(path)); err != nil {
		return fmt.Errorf("failed to download media: %w", err)
	}

	if err := replaygain(cp, config.Command.Replaygain, path); err != nil {
		return fmt.Errorf("failed to tag replaygain information: %w", err)
	}

	return nil
}
