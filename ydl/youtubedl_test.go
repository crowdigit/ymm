package ydl_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/crowdigit/ymm/mock"
	"github.com/crowdigit/ymm/ydl"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type YoutubeDLTestSuite struct {
	suite.Suite

	mockCtrl            *gomock.Controller
	mockCommandProvider *mock.MockCommandProvider
}

func (s *YoutubeDLTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockCommandProvider = mock.NewMockCommandProvider(s.mockCtrl)
}

func (s *YoutubeDLTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *YoutubeDLTestSuite) TestPlaylistMetadata() {
	// TODO
}

func (s *YoutubeDLTestSuite) TestVideoMetadata() {
	testOutputFile, err := os.Open("./test.json")
	s.Nil(err)

	testOutput, err := io.ReadAll(testOutputFile)
	s.Nil(err)

	url := "http://asdf.com/some/url"
	s.mockCommandProvider.EXPECT().
		NewCommand("youtube-dl", "--dump-json", url).
		DoAndReturn(func(name string, args ...string) ydl.Command {
			command := mock.NewMockCommand(s.mockCtrl)
			command.EXPECT().Start().Times(1)
			command.EXPECT().StderrPipe().
				Return(io.NopCloser(bytes.NewReader(nil)), nil).
				Times(1)
			command.EXPECT().StdoutPipe().
				Return(io.NopCloser(bytes.NewReader(testOutput)), nil).
				Times(1)
			command.EXPECT().Wait().Times(1)
			return command
		}).
		Times(1)

	ydl := ydl.NewYoutubeDLImpl(s.mockCommandProvider)
	got, err := ydl.VideoMetadata(url)
	s.Nil(err)

	expect := commonTestMetadata
	s.Equal(expect, got)
}

func (s *YoutubeDLTestSuite) TestDownload() {
	// TODO
}

func TestYoutubeDLTestSuite(t *testing.T) {
	suite.Run(t, new(YoutubeDLTestSuite))
}
