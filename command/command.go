package command

import (
	"fmt"
	"io"
	"os/exec"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=../mock/mock_command.go -package=mock github.com/crowdigit/ymm/command Command
type Command interface {
	SetDir(string)
	Start() error
	StderrPipe() (io.ReadCloser, error)
	StdoutPipe() (io.ReadCloser, error)
	StdinPipe() (io.WriteCloser, error)
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

func (c commandImpl) StdinPipe() (io.WriteCloser, error) {
	return c.cmd.StdinPipe()
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

func handleStream(
	logger *zap.SugaredLogger,
	chStdoutIn <-chan []byte,
	chStdoutOut chan<- []byte,
	chClose chan struct{},
	chErrIn <-chan error,
	chErrOut chan<- error) {
	result := make([]byte, 0, 8192)
	var err error
loop:
	for {
		select {
		case err = <-chErrIn:
			logger.Errorf("reading from command stdout returned an error: %s", err)
			break loop
		case msg := <-chStdoutIn:
			result = append(result, msg...)
		case <-chClose:
			break loop
		}
	}
	chStdoutOut <- result
	chErrOut <- err
}

//Run runs command and returns stdout synchronously, dropping stderr
func Run(logger *zap.SugaredLogger, command Command) ([]byte, error) {
	stdout, err := command.StdoutPipe()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stdout pipe")
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	chStdout := make(chan []byte)
	chStdoutResult := make(chan []byte)
	chErr := make(chan error)
	chErrResult := make(chan error)
	go ReadStream(&wg, stdout, chStdout, chErr)

	if err := command.Start(); err != nil {
		return nil, errors.Wrap(err, "failed to start command")
	}

	chClose := make(chan struct{})
	go handleStream(logger, chStdout, chStdoutResult, chClose, chErr, chErrResult)

	wg.Wait()
	close(chClose)

	result := <-chStdoutResult
	err = <-chErrResult
	if err != nil {
		return nil, errors.Wrap(err, "stream handler returned error")
	}

	status, err := command.Wait()
	if err != nil {
		return nil, errors.Wrap(err, "failed to wait for command")
	}

	if status != 0 {
		return nil, fmt.Errorf("command exited with %d", status)
	}

	return result, nil
}
