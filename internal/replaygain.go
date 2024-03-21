package internal

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/crowdigit/exec"
)

func replaygain(cp exec.CommandProvider, replaygain ConfigExec, path string) error {
	ctx, kill := context.WithCancel(context.Background())
	defer kill()
	ctx, unregister := signal.NotifyContext(ctx, os.Interrupt)
	defer unregister()

	replaygain = replaygain.ReplacePlaceholder("<path>", path)

	cmd := cp.CommandContext(ctx, exec.CommandOpts{
		Path:   replaygain.Path,
		Args:   replaygain.Args,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run replaygain command: %w", err)
	}

	return nil
}
