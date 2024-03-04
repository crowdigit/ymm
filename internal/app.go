package internal

import (
	"bytes"
	"context"
	"fmt"
	"os"
	oexec "os/exec"
	"os/signal"

	"github.com/crowdigit/ymm/pkg/exec"
)

type ExecConfig struct {
	Path string
	Args []string
}

func DownloadSingle(ci exec.CommandProvider, jqConf, ytConf ExecConfig, url string) error {
	ctx, kill := context.WithCancel(context.Background())
	defer kill()
	ctx, unregister := signal.NotifyContext(ctx, os.Interrupt)
	defer unregister()

	ytOpts := exec.CommandOpts{
		Path: ytConf.Path,
		Args: append(ytConf.Args, url),
	}
	ytCmd := ci.CommandContext(ctx, ytOpts)
	ytOut, err := ytCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get yt command stdout pipe: %w", err)
	}

	if err := ytCmd.Start(); err != nil {
		if err, ok := err.(*oexec.ExitError); ok {
			fmt.Fprintln(os.Stderr, string(err.Stderr))
		}
		return fmt.Errorf("failed to start yt command: %w", err)
	}

	var ytStdout bytes.Buffer
	chYtStdoutRead := make(chan error)
	go read("yt", ytOut, &ytStdout, chYtStdoutRead)

	var pipeErr error
	select {
	case <-ctx.Done():
		return fmt.Errorf("interrupted")
	case pipeErr = <-chYtStdoutRead:
	}

	if err := ytCmd.Wait(); err != nil {
		if err, ok := err.(*oexec.ExitError); ok {
			fmt.Fprintf(os.Stderr, "yt returned with status %d\n", err.ExitCode())
			fmt.Fprintln(os.Stderr, string(err.Stderr))
			return fmt.Errorf("yt command returned with status code %d", err.ExitCode())
		}
		return fmt.Errorf("failed to execute yt command: %w", err)
	}

	if pipeErr != nil {
		return fmt.Errorf("failed to read from yt command stdout: %w", pipeErr)
	}

	n := min(30, ytStdout.Len())
	fmt.Fprintf(os.Stdout, "read: %s\n", ytStdout.String()[:n])

	return nil
}
