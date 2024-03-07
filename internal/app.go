package internal

import (
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

	pipeline, err := NewPipeline(ctx, cp, []PipeSpec{
		{CmdOpt: exec.CommandOpts{Path: ytConf.Path, Args: append(ytConf.Args, url)}},
		{CmdOpt: exec.CommandOpts{Path: jqConf.Path, Args: jqConf.Args}},
	})
	if err != nil {
		return fmt.Errorf("failed to initialize pipeline: %w", err)
	} else if err := pipeline.Start(); err != nil {
		kill()
		_ = pipeline.Wait()
		return fmt.Errorf("failed to start pipeline: %w", err)
	}

	chErr := make(chan error)
	go func() {
		for {
			subBuffer := make([]byte, 1024)
			read, err := pipeline.output.Read(subBuffer)
			if read > 0 {
				fmt.Println(string(subBuffer[:read]))
			}
			if errors.Is(err, io.EOF) {
				fmt.Println("EOF")
				close(chErr)
				return
			} else if err != nil {
				fmt.Printf("error: %s\n", err)
				chErr <- err
				return
			}
		}
	}()

	<-chErr

	_ = pipeline.Wait()

	return nil
}
