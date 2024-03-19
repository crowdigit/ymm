package internal

import (
	"fmt"

	"github.com/crowdigit/exec"
)

type ExecConfig struct {
	Path string
	Args []string
}

func DownloadSingle(cp exec.CommandProvider, ytConf, jqConf ExecConfig, url string) error {
	for i := range ytConf.Args {
		if ytConf.Args[i] == "<url>" {
			ytConf.Args[i] = url
		}
	}

	jqMetadata, err := fetchMetadata(cp, ytConf, jqConf, url)
	if err != nil {
		return fmt.Errorf("failed to fetch media metadata: %w", err)
	}

	if err := validateFormats(jqMetadata); err != nil {
		return fmt.Errorf("failed to validate available media formats: %w", err)
	}

	fmt.Println("perfectly valid video!")

	return nil
}
