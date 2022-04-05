package jq

import (
	"errors"

	"github.com/crowdigit/ymm/command"
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
	return nil, errors.New("JqImpl.Slurp is not implemented")
}
