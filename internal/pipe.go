package internal

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/crowdigit/ymm/pkg/exec"
)

type StreamType int

const (
	Stdout StreamType = iota
	Stderr StreamType = iota
	Null   StreamType = iota
)

func (t StreamType) String() string {
	switch t {
	case Stdout:
		return "stdout"
	case Stderr:
		return "stderr"
	}
	return "invalid stream type"
}

type PipeSpec struct {
	CmdOpt exec.CommandOpts
	Next   StreamType
	Other  io.Writer
}

type Pipeline struct {
	cmds   []exec.Command
	output io.Reader
}

func NewPipeline(
	ctx context.Context,
	cp exec.CommandProvider,
	pipeSpecs []PipeSpec,
) (Pipeline, error) {
	if len(pipeSpecs) == 0 {
		return Pipeline{}, nil
	}

	prev := pipeSpecs[0].CmdOpt.Stdin
	var err error

	cmds := make([]exec.Command, 0, len(pipeSpecs))
	for i, pipeSpec := range pipeSpecs {
		pipeSpec.CmdOpt.Stdin = prev
		cmd := cp.CommandContext(ctx, pipeSpec.CmdOpt)
		cmds = append(cmds, cmd)
		switch pipeSpec.Next {
		case Stdout:
			prev, err = cmd.StdoutPipe()
		case Stderr:
			prev, err = cmd.StderrPipe()
		case Null:
			if i != len(pipeSpecs)-1 {
				return Pipeline{}, errors.New("attempted to pipe to null for non-terminal pipe")
			}
			prev = nil
		}
		if err != nil {
			return Pipeline{}, fmt.Errorf("failed to get standard stream %d: %w", i, err)
		}
	}
	return Pipeline{cmds, prev}, nil
}

func (p Pipeline) Start() error {
	var errs []error
	for _, cmd := range p.cmds {
		if err := cmd.Start(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (p Pipeline) Wait() error {
	var errs []error
	for _, cmd := range p.cmds {
		if err := cmd.Wait(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
