package loudness

import (
	"github.com/crowdigit/ymm/command"
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

func (s LoudnessScannerImpl) Tag(path string) error {
	panic("not implemented") // TODO: Implement
}
