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

	chErr       chan error
	startWaiter sync.Once
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
		cmd := cp.CommandContext(ctx, pipeSpec.CmdOpt)
		cmds = append(cmds, cmd)
		switch pipeSpec.Next {
		case Stdout:
			prev, err = cmd.StdoutPipe()
		case Stderr:
			prev, err = cmd.StderrPipe()
		case Null:
			if i != len(pipeSpecs)-1 {
				return nil, errors.New("attempted to pipe to null for non-terminal pipe")
			}
			prev = nil
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get standard stream %d: %w", i, err)
		}
	}
	return &Pipeline{
		cmds:        cmds,
		output:      prev,
		chErr:       make(chan error),
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
	chErr := make(chan error)
	for _, cmd := range p.cmds {
		go func(cmd exec.Command) {
			chErr <- cmd.Wait()
		}(cmd)
	}
	for i := 0; i < len(p.cmds); i += 1 {
		if err := <-chErr; err != nil {
			p.chErr <- err
		}
	}
	close(p.chErr)
}

func (p *Pipeline) Wait() <-chan error {
	p.startWaiter.Do(func() { go p.wait() })
	return p.chErr
}
