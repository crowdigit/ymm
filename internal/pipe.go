package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

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

	chErr       chan PipelineError
	startWaiter sync.Once
}

func nextStdout(
	ctx context.Context,
	cp exec.CommandProvider,
	index int,
	pipeSpec *PipeSpec,
) (exec.Command, io.ReadCloser, error) {
	pipeSpec.CmdOpt.Stderr = pipeSpec.Other
	cmd := cp.CommandContext(ctx, pipeSpec.CmdOpt)
	prev, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get standard stream %d: %w", index, err)
	}
	return cmd, prev, nil
}

func nextStderr(
	ctx context.Context,
	cp exec.CommandProvider,
	index int,
	pipeSpec *PipeSpec,
) (exec.Command, io.ReadCloser, error) {
	pipeSpec.CmdOpt.Stdout = pipeSpec.Other
	cmd := cp.CommandContext(ctx, pipeSpec.CmdOpt)
	prev, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get standard stream %d: %w", index, err)
	}
	return cmd, prev, nil
}

func nextNull(
	ctx context.Context,
	cp exec.CommandProvider,
	pipeSpecsLen int,
	index int,
	pipeSpec *PipeSpec,
) (exec.Command, io.ReadCloser, error) {
	cmd := cp.CommandContext(ctx, pipeSpec.CmdOpt)
	if pipeSpecsLen-1 != index {
		return nil, nil, errors.New("attempted to pipe to null for non-terminal pipe")
	}
	return cmd, nil, nil
}

func NewPipeline(
	ctx context.Context,
	cp exec.CommandProvider,
	pipeSpecs []PipeSpec,
) (*Pipeline, error) {
	if len(pipeSpecs) == 0 {
		return nil, errors.New("len(pipeSpecs) == 0")
	}

	prev := pipeSpecs[0].CmdOpt.Stdin
	var err error
	cmds := make([]exec.Command, 0, len(pipeSpecs))
	for i, pipeSpec := range pipeSpecs {
		pipeSpec.CmdOpt.Stdin = prev
		var cmd exec.Command
		switch pipeSpec.Next {
		case Stdout:
			cmd, prev, err = nextStdout(ctx, cp, i, &pipeSpec)
		case Stderr:
			cmd, prev, err = nextStderr(ctx, cp, i, &pipeSpec)
		case Null:
			cmd, prev, err = nextNull(ctx, cp, len(pipeSpecs), i, &pipeSpec)
		}
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, cmd)
	}
	return &Pipeline{
		cmds:        cmds,
		output:      prev,
		chErr:       make(chan PipelineError),
		startWaiter: sync.Once{},
	}, nil
}

func (p *Pipeline) Start() error {
	var errs []error
	for _, cmd := range p.cmds {
		if err := cmd.Start(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (p *Pipeline) wait() {
	chErr := make(chan PipelineError)
	for index, cmd := range p.cmds {
		go func(index int, cmd exec.Command) {
			chErr <- PipelineError{cmd.Wait(), index}
		}(index, cmd)
	}
	for i := 0; i < len(p.cmds); i += 1 {
		if err := <-chErr; err.err != nil {
			p.chErr <- err
		}
	}
	close(p.chErr)
}

// Wait returns channel which will send errors caused by command. It closes
// if every process in pipeline has exited. Closing without sending any error
// means all processes have exited normally. It is safe to call Wait multiple
// times.
//
// Errors sent over returned channel is guranteed to be type of [PipelineError]
// which wraps original error.
func (p *Pipeline) Wait() <-chan PipelineError {
	p.startWaiter.Do(func() { go p.wait() })
	return p.chErr
}
