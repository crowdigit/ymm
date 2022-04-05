package jq

import (
	"github.com/crowdigit/ymm/command"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=../mock/mock_jq.go -package=mock github.com/crowdigit/ymm/jq Jq
type Jq interface {
	Slurp([]byte) ([]byte, error)
}

type JqImpl struct {
	logger          *zap.SugaredLogger
	commandProvider command.CommandProvider
}

func NewJq(logger *zap.SugaredLogger, commandProvider command.CommandProvider) Jq {
	return &JqImpl{
		logger:          logger,
		commandProvider: commandProvider,
	}
}

func (jq *JqImpl) Slurp([]byte) ([]byte, error) {
	cmd := jq.commandProvider.NewCommand("jq", "--slurp", ".")
	stdout, err := command.Run(jq.logger, cmd)
	return stdout, errors.Wrap(err, "failed to run command")
}
