package ydl_test

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/crowdigit/ymm/command"
	"github.com/crowdigit/ymm/mock"
	"github.com/crowdigit/ymm/ydl"
	"github.com/golang/mock/gomock"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
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

func (s *YoutubeDLTestSuite) TestVideoMetadata() {
	testOutputFile, err := os.Open("./test.json")
	s.Nil(err)

	testOutput, err := io.ReadAll(testOutputFile)
	s.Nil(err)

	url := "http://asdf.com/some/url"

	s.mockCommandProvider.EXPECT().
		NewCommand(
			"yt-dlp",
			"-o",
			"%(title)s-%(id)s.%(ext)s",
			"--dump-json",
			url).
		DoAndReturn(func(name string, args ...string) command.Command {
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

	youtubeDl := ydl.NewYoutubeDLImpl(zap.NewNop().Sugar(), s.mockCommandProvider, "")
	result, err := youtubeDl.VideoMetadata(url)
	s.Nil(err)

	got := ydl.VideoMetadata{}
	s.Nil(jsoniter.Unmarshal(result, &got))

	expected := commonTestMetadata
	s.Equal(expected, got)
}

func (s *YoutubeDLTestSuite) TestDownload() {
	metadata := ydl.VideoMetadata{
		ID:      "456",
		Title:   "some title",
		Formats: []ydl.Format{{FormatID: "251"}},
		Thumbnails: []ydl.Thumbnail{
			{
				URL:    "http://asdf.com/thunbnail/456.png",
				ID:     "1",
				Height: 1080,
				Width:  1920,
			},
		},
		Description: "some description",
		Uploader:    "some uploader",
		UploaderID:  "123",
		UploaderURL: "http://asdf.com/some/uploader/123",
		UploadDate:  ydl.NewJSONTime(time.Date(2022, time.March, 1, 1, 0, 0, 0, time.UTC)),
		WebpageURL:  "http://asdf.com/some/url/456",
		Filename:    "somefilename.mp3",
		Duration:    0,
	}
	downloadDir := "/asdf"

	s.mockCommandProvider.EXPECT().
		NewCommand(
			"yt-dlp",
			"-o",
			"%(title)s-%(id)s.%(ext)s",
			"--format", "251",
			"--extract-audio",
			"--audio-format", "mp3",
			"--audio-quality", "0",
			metadata.WebpageURL,
		).
		DoAndReturn(func(name string, args ...string) command.Command {
			command := mock.NewMockCommand(s.mockCtrl)
			command.EXPECT().SetDir(downloadDir)
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

	youtubeDl := ydl.NewYoutubeDLImpl(zap.NewNop().Sugar(), s.mockCommandProvider, "")
	got, err := youtubeDl.Download(downloadDir, metadata)
	s.Nil(err)

	expected := ydl.DownloadResult{}
	s.Equal(expected, got)
}

func TestYoutubeDLTestSuite(t *testing.T) {
	suite.Run(t, new(YoutubeDLTestSuite))
}
