package app_test

import (
	"testing"

	"github.com/crowdigit/ymm/app"
	"github.com/crowdigit/ymm/mock"
	"github.com/crowdigit/ymm/ydl"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type AppTestSuite struct {
	suite.Suite

	mockCtrl *gomock.Controller
	mockYdl  *mock.MockYoutubeDL
}

func (s *AppTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockYdl = mock.NewMockYoutubeDL(s.mockCtrl)
}

func (s *AppTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *AppTestSuite) TestDownloadSingle() {
	url := "http://asdf.com/some/url"
	metadata := ydl.VideoMetadata{}
	result := ydl.DownloadResult{}

	s.mockYdl.EXPECT().VideoMetadata(url).
		DoAndReturn(func(url string) (ydl.VideoMetadata, error) {
			return metadata, nil
		}).
		Times(1)
	s.mockYdl.EXPECT().Download(metadata).
		DoAndReturn(func(metadata ydl.VideoMetadata) (ydl.DownloadResult, error) {
			return result, nil
		}).
		Times(1)

	app := app.NewApplicationImpl(s.mockYdl)
	s.Nil(app.DownloadSingle(url))
}

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}
