package internal

import (
	"fmt"

	"github.com/crowdigit/exec"
)

func DownloadSingle(cp exec.CommandProvider, config Config, url string) error {
	config.Command.Metadata.Youtube.ReplacePlaceholder("<url>", url)
	config.Command.Download.Youtube.ReplacePlaceholder("<url>", url)

	jqMetadata, err := fetchMetadata(
		cp,
		config.Command.Metadata.Youtube,
		config.Command.Metadata.JSON,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch media metadata: %w", err)
	}

	format, err := validateFormats(jqMetadata)
	if err != nil {
		return fmt.Errorf("failed to validate available media formats: %w", err)
	}
	config.Command.Download.Youtube.ReplacePlaceholder("<format>", format)

	if err := downloadVideo(cp, config.Command.Download.Youtube); err != nil {
		return fmt.Errorf("failed to download media: %w", err)
	}

	// TODO persist metadata into DB

	return nil
}
