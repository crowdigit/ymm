package internal

import (
	"fmt"

	"github.com/crowdigit/exec"
)

type ExecConfig struct {
	Path string
	Args []string
}

func replacePlaceholder(args []string, placeholder, url string) {
	for i := range args {
		if args[i] == placeholder {
			args[i] = url
		}
	}
}

func DownloadSingle(
	cp exec.CommandProvider,
	cmdYtMetadata, cmdJq, cmdYtDownload ExecConfig,
	url string,
) error {
	replacePlaceholder(cmdYtMetadata.Args, "<url>", url)
	replacePlaceholder(cmdYtDownload.Args, "<url>", url)

	jqMetadata, err := fetchMetadata(cp, cmdYtMetadata, cmdJq)
	if err != nil {
		return fmt.Errorf("failed to fetch media metadata: %w", err)
	}

	format, err := validateFormats(jqMetadata)
	if err != nil {
		return fmt.Errorf("failed to validate available media formats: %w", err)
	}

	replacePlaceholder(cmdYtDownload.Args, "<format>", format)
	if err := downloadVideo(cp, cmdYtDownload); err != nil {
		return fmt.Errorf("failed to download media: %w", err)
	}

	// TODO persist metadata into DB

	return nil
}
