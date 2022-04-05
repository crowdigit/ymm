package jq_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/crowdigit/ymm/command"
	"github.com/crowdigit/ymm/jq"
	"github.com/crowdigit/ymm/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type JqTestSuite struct {
	suite.Suite

	mockCtrl            *gomock.Controller
	mockCommandProvider *mock.MockCommandProvider
}

func (s *JqTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockCommandProvider = mock.NewMockCommandProvider(s.mockCtrl)
}

func (s *JqTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *JqTestSuite) TestSlurp() {
	wackyInput := []byte(`{
		"lol": 1234
	}
	{
		"no one": "tell me what to do"
	}`)
	noneWackyOutput := []byte(`{
		"lol": 1234
	},
	{
		"no one": "tell me what to do"
	}`)

	s.mockCommandProvider.EXPECT().
		NewCommand("jq", "--slurp", ".").
		DoAndReturn(func(name string, args ...string) command.Command {
			command := mock.NewMockCommand(s.mockCtrl)
			command.EXPECT().Start().Times(1)
			command.EXPECT().StderrPipe().
				Return(io.NopCloser(bytes.NewReader(nil)), nil).
				Times(1)
			command.EXPECT().StdoutPipe().
				Return(io.NopCloser(bytes.NewReader(noneWackyOutput)), nil).
				Times(1)
			command.EXPECT().Wait().Times(1)
			return command
		}).
		Times(1)

	jq := jq.NewJq(zap.NewNop().Sugar(), s.mockCommandProvider)
	s.Nil(jq.Slurp(wackyInput))
}

func TestJqTestSuite(t *testing.T) {
	suite.Run(t, new(JqTestSuite))
}
