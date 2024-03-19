package internal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/crowdigit/exec"
)

type ExecConfig struct {
	Path string
	Args []string
}

func DownloadSingle(cp exec.CommandProvider, jqConf, ytConf ExecConfig, url string) error {
	ctx, kill := context.WithCancel(context.Background())
	defer kill()
	ctx, unregister := signal.NotifyContext(ctx, os.Interrupt)
	defer unregister()

	buffers := []*bytes.Buffer{{}, {}}
	pipespecs := []exec.PipeSpec{
		{
			CmdOpt: exec.CommandOpts{Path: ytConf.Path, Args: append(ytConf.Args, url)},
			Next:   exec.Stdout,
			Other:  buffers[0],
		},
		{
			CmdOpt: exec.CommandOpts{Path: jqConf.Path, Args: jqConf.Args},
			Next:   exec.Stdout,
			Other:  buffers[1],
		},
	}

	pipeline, err := exec.NewPipeline(ctx, cp, pipespecs)
	if err != nil {
		return fmt.Errorf("failed to initialize pipeline: %w", err)
	}

	err = pipeline.Start()
	defer pipeline.Cancel(kill)
	if err != nil {
		return fmt.Errorf("failed to start pipeline: %w", err)
	}

	chPipeErr := make(chan error)
	go func() {
		for {
			subBuffer := make([]byte, 1024)
			read, err := pipeline.Output().Read(subBuffer)
			if read > 0 {
				fmt.Println(string(subBuffer[:read]))
			}
			if errors.Is(err, io.EOF) {
				close(chPipeErr)
				return
			} else if err != nil {
				chPipeErr <- err
				return
			}
		}
	}()
	if err := <-chPipeErr; err != nil {
		return fmt.Errorf("failed to operate on pipe: %w", err)
	}

	if errs := pipeline.Any(kill); errs != nil {
		fmt.Fprintln(os.Stderr, errs)
		for i, err := range errs {
			fmt.Println(err)
			if buffers[errs[i].Index].Len() > 0 {
				fmt.Println(buffers[i].String())
			}
		}
		return errs
	}

	return nil
}
