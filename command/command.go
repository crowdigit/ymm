package command

import (
	"io"
	"os/exec"
	"sync"

	"github.com/pkg/errors"
)

//go:generate mockgen -destination=../mock/mock_command.go -package=mock github.com/crowdigit/ymm/command Command
type Command interface {
	SetDir(string)
	Start() error
	StderrPipe() (io.ReadCloser, error)
	StdoutPipe() (io.ReadCloser, error)
	Wait() (int, error)
}

//go:generate mockgen -destination=../mock/mock_command_provider.go -package=mock github.com/crowdigit/ymm/command CommandProvider
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

func ReadStream(wg *sync.WaitGroup, reader io.ReadCloser, chOut chan<- []byte, chErr chan<- error) {
	defer wg.Done()
	readBuffer := make([]byte, 1024)
	for {
		read, err := reader.Read(readBuffer)
		if read > 0 {
			sendBuffer := make([]byte, read)
			copy(sendBuffer, readBuffer)
			chOut <- sendBuffer
		} else if err == io.EOF {
			break
		} else if err != nil {
			chErr <- errors.Wrap(err, "failed to read from reader stream")
			break
		}
	}
}
