package loudness

import (
	"fmt"
	"sync"

	"github.com/crowdigit/ymm/command"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=../mock/mock_loudness.go -package=mock github.com/crowdigit/ymm/loudness LoudnessScanner
type LoudnessScanner interface {
	Tag(path string) error
}

type LoudnessScannerImpl struct {
	logger          *zap.SugaredLogger
	commandProvider command.CommandProvider
}

func NewLoudnessScanner(logger *zap.SugaredLogger, commandProvider command.CommandProvider) LoudnessScanner {
	return LoudnessScannerImpl{
		logger:          logger,
		commandProvider: commandProvider,
	}
}

func handleTagStream(logger *zap.SugaredLogger, chStdout <-chan []byte, chStderr <-chan []byte, chClose chan struct{}, chErr <-chan error) {
loop:
	for {
		select {
		case err := <-chErr:
			logger.Errorf("reading from loudness scanner stream returned an error: %s", err)
			break loop
		case msg := <-chStderr:
			logger.Warnf(string(msg))
		case msg := <-chStdout:
			logger.Infof(string(msg))
		case <-chClose:
			break loop
		}
	}
}

func (s LoudnessScannerImpl) Tag(path string) error {
	cmd := s.commandProvider.NewCommand("loudness", "tag", "--track", path)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return errors.Wrap(err, "failed to get stderr pipe")
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "failed to get stdout pipe")
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	chStderr := make(chan []byte)
	chStdout := make(chan []byte)
	chErr := make(chan error)
	go command.ReadStream(&wg, stderr, chStderr, chErr)
	go command.ReadStream(&wg, stdout, chStdout, chErr)

	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, "failed to start loudness scanner command")
	}

	chClose := make(chan struct{})
	go handleTagStream(s.logger, chStdout, chStderr, chClose, chErr)

	wg.Wait()

	close(chClose)

	status, err := cmd.Wait()
	if err != nil {
		return errors.Wrap(err, "failed to wait for loudness scanner command")
	}

	if status != 0 {
		return fmt.Errorf("loudness scanner command exited with %d", status)
	}

	return nil
}
