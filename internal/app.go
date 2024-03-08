package internal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/crowdigit/ymm/pkg/exec"
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

	buffers := []*bytes.Buffer{{}, {}, {}}
	pipeline, err := NewPipeline(ctx, cp, []PipeSpec{
		{
			CmdOpt: exec.CommandOpts{Path: ytConf.Path, Args: append(ytConf.Args, url)},
			Next:   Stdout,
			Other:  buffers[0],
		},
		{
			CmdOpt: exec.CommandOpts{Path: jqConf.Path, Args: jqConf.Args},
			Next:   Stdout,
			Other:  buffers[1],
		},
	})
	if err != nil {
		return fmt.Errorf("failed to initialize pipeline: %w", err)
	}

	err = pipeline.Start()
	defer func() {
		kill()
		for range pipeline.Wait() {
		}
	}()
	if err != nil {
		return fmt.Errorf("failed to start pipeline: %w", err)
	}

	chPipeErr := make(chan error)
	go func() {
		for {
			subBuffer := make([]byte, 1024)
			read, err := pipeline.output.Read(subBuffer)
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

	errs := make([]error, 0, len(pipeline.cmds))
	for err := range pipeline.Wait() {
		kill()
		errs = append(errs, err)
	}

	fmt.Printf("--- buffer[0]\n%s", buffers[0].String())
	fmt.Println("---")
	fmt.Printf("--- buffer[1]\n%s", buffers[1].String())
	fmt.Println("---")

	return errors.Join(errs...)
}
