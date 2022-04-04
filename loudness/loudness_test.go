package loudness_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/crowdigit/ymm/command"
	"github.com/crowdigit/ymm/loudness"
	"github.com/crowdigit/ymm/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type LoudnessScannerTestSuite struct {
	suite.Suite

	mockCtrl            *gomock.Controller
	mockCommandProvider *mock.MockCommandProvider
}

func (s *LoudnessScannerTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockCommandProvider = mock.NewMockCommandProvider(s.mockCtrl)
}

func (s *LoudnessScannerTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *LoudnessScannerTestSuite) TestTag() {
	path := "some filepath"

	s.mockCommandProvider.EXPECT().
		NewCommand("loudness", "tag", "--track", path).
		DoAndReturn(func(name string, args ...string) command.Command {
			command := mock.NewMockCommand(s.mockCtrl)
			command.EXPECT().Start().Times(1)
			command.EXPECT().StderrPipe().
				Return(io.NopCloser(bytes.NewReader(nil)), nil).
				Times(1)
			command.EXPECT().StdoutPipe().
				Return(io.NopCloser(bytes.NewReader(nil)), nil).
				Times(1)
			command.EXPECT().Wait().Times(1)
			return command
		}).
		Times(1)

	loudnessScanner := loudness.NewLoudnessScanner(zap.NewNop().Sugar(), s.mockCommandProvider)
	s.Nil(loudnessScanner.Tag(path))
}

func TestLoudnessScanner(t *testing.T) {
	suite.Run(t, new(LoudnessScannerTestSuite))
}
