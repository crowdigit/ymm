package internal

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/crowdigit/exec"
)

func downloadVideo(
	cp exec.CommandProvider,
	youtube ConfigExec,
	url string,
	format string,
	dir string,
) error {
	ctx, kill := context.WithCancel(context.Background())
	defer kill()
	ctx, unregister := signal.NotifyContext(ctx, os.Interrupt)
	defer unregister()

	youtube = youtube.ReplacePlaceholder("<url>", url).ReplacePlaceholder("<format>", format)

	cmd := cp.CommandContext(ctx, exec.CommandOpts{
		Path:   youtube.Path,
		Args:   youtube.Args,
		Dir:    dir,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start download command: %w", err)
	}

	return nil
}
