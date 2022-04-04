package ydl

import (
	"io"
	"os/exec"
)

//go:generate mockgen -destination=../mock/mock_command.go -package=mock github.com/crowdigit/ymm/ydl Command
type Command interface {
	SetDir(string)
	Start() error
	StderrPipe() (io.ReadCloser, error)
	StdoutPipe() (io.ReadCloser, error)
	Wait() (int, error)
}

//go:generate mockgen -destination=../mock/mock_command_provider.go -package=mock github.com/crowdigit/ymm/ydl CommandProvider
type CommandProvider interface {
	NewCommand(name string, arg ...string) Command
}

type commandImpl struct {
	cmd *exec.Cmd
}

func (c commandImpl) SetDir(dir string) {
	c.cmd.Dir = dir
}

func (c commandImpl) Start() error {
	return c.cmd.Start()
}

func (c commandImpl) StderrPipe() (io.ReadCloser, error) {
	return c.cmd.StderrPipe()
}

func (c commandImpl) StdoutPipe() (io.ReadCloser, error) {
	return c.cmd.StdoutPipe()
}

func (c commandImpl) Wait() (int, error) {
	if err := c.cmd.Wait(); err != nil {
		return 0, err
	}
	code := c.cmd.ProcessState.ExitCode()
	return code, nil
}

type commandProviderImpl struct{}

func (cp commandProviderImpl) NewCommand(name string, arg ...string) Command {
	return commandImpl{
		cmd: exec.Command(name, arg...),
	}
}

func NewCommandProviderImpl() commandProviderImpl {
	return commandProviderImpl{}
}
