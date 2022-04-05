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

func (jq *JqImpl) Slurp(input []byte) ([]byte, error) {
	jq.logger.Info("entered JqImpl.Slurp")
	defer jq.logger.Info("exiting JqImpl.Slurp")
	cmd := jq.commandProvider.NewCommand("jq", "--slurp", ".")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stdin of jq command")
	}

	go func() {
		defer stdin.Close()
		stdin.Write(input)
	}()

	stdout, err := command.Run(jq.logger, cmd)
	jq.logger.Infof("jq result: %s", string(stdout))
	return stdout, errors.Wrap(err, "failed to run command")
}
