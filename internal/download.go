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
	cmdYtDownload ExecConfig,
) error {
	ctx, kill := context.WithCancel(context.Background())
	defer kill()
	ctx, unregister := signal.NotifyContext(ctx, os.Interrupt)
	defer unregister()

	cmd := cp.CommandContext(ctx, exec.CommandOpts{
		Path:   cmdYtDownload.Path,
		Args:   cmdYtDownload.Args,
		Dir:    "",
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start download command: %w", err)
	}

	return nil
}
