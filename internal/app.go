package internal

import (
	"fmt"

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

	if err := downloadVideo(cp, config.Command.Download.Youtube, url, format); err != nil {
		return fmt.Errorf("failed to download media: %w", err)
	}

	// TODO persist metadata into DB

	return nil
}
