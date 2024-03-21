package internal

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/crowdigit/exec"
)

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

	if err := downloadVideo(cp, config.Command.Download.Youtube, url, format); err != nil {
		return fmt.Errorf("failed to download media: %w", err)
	}

	path := strings.TrimSuffix(
		jqMetadata[0].Filename,
		filepath.Ext(jqMetadata[0].Filename),
	) + ".mp3"

	if err := replaygain(cp, config.Command.Replaygain, path); err != nil {
		return fmt.Errorf("failed to tag replaygain information: %w", err)
	}

	return nil
}
