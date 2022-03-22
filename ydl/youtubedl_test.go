package ydl_test

import (
	"testing"

	"github.com/crowdigit/ymm/mock"
	"github.com/crowdigit/ymm/ydl"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type YoutubeDLTestSuite struct {
	suite.Suite

	mockCtrl    *gomock.Controller
	mockCommand *mock.MockCommand
}

func (s *YoutubeDLTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockCommand = mock.NewMockCommand(s.mockCtrl)
}

func (s *YoutubeDLTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *YoutubeDLTestSuite) TestPlaylistMetadata() {
}

func (s *YoutubeDLTestSuite) TestVideoMetadata() {
	url := "http://asdf.com/some/url"
	s.mockCommand.EXPECT().Start("youtube-dl", "--dump-json", url).Times(1)
	s.mockCommand.EXPECT().Wait().Times(1)

	ydl.NewYoutubeDLImpl(s.mockCommand)
}

func (s *YoutubeDLTestSuite) TestDownload() {
}

func TestYoutubeDLTestSuite(t *testing.T) {
	suite.Run(t, new(YoutubeDLTestSuite))
}
