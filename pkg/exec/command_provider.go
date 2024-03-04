package sysutil

import (
	"context"
	io "io"
	"os/exec"
)

//go:generate mockgen -source=$GOFILE -package $GOPACKAGE -destination mock_$GOFILE

type CommandOpts struct {
	Path   string
	Args   []string
	Dir    string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type CommandProvider interface {
	Command(opt CommandOpts) Command
	CommandContext(ctx context.Context, opt CommandOpts) Command
}

type Command interface {
	Output() ([]byte, error)

	Start() error
	Wait() error
	Run() error

	StdinPipe() (io.WriteCloser, error)
	StdoutPipe() (io.ReadCloser, error)
	StderrPipe() (io.ReadCloser, error)

	PID() int
}

type commandImpl struct {
	*exec.Cmd
}

func (c commandImpl) PID() int {
	return c.Process.Pid
}

type commandProvierImpl struct{}

func (i commandProvierImpl) Command(opts CommandOpts) Command {
	cmd := exec.Command(opts.Path, opts.Args...)
	cmd.Dir = opts.Dir
	cmd.Stdin = opts.Stdin
	cmd.Stdout = opts.Stdout
	cmd.Stderr = opts.Stderr
	return commandImpl{cmd}
}

func (i commandProvierImpl) CommandContext(ctx context.Context, opt CommandOpts) Command {
	cmd := exec.CommandContext(ctx, opt.Path, opt.Args...)
	cmd.Dir = opt.Dir
	cmd.Stdin = opt.Stdin
	cmd.Stdout = opt.Stdout
	cmd.Stderr = opt.Stderr
	return commandImpl{cmd}
}

func NewCommandProvider() commandProvierImpl {
	return commandProvierImpl{}
}
