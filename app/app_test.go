package app_test

import (
	"testing"
	"time"

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
	metadataUploadDate := ydl.JSONTime(time.Date(2022, time.February, 27, 0, 0, 0, 0, time.UTC))
	metadata := ydl.VideoMetadata{
		ID:      "some_id",
		Title:   "some title",
		Formats: []ydl.Format{{FormatID: "251"}},
		Thumbnails: []ydl.Thumbnail{{
			URL:    "some url",
			ID:     "1",
			Height: 1080,
			Width:  1920,
		}},
		Description: "some description",
		Uploader:    "some uploader",
		UploaderID:  "some_uploader",
		UploaderURL: "some_uploader_url",
		UploadDate:  &metadataUploadDate,
		Filename:    "some_filename",
		Duration:    123,
	}
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

func (s *AppTestSuite) TestDownloadPlaylist() {
	url := "http://asdf.com/some/url"

	metadata := []ydl.VideoMetadata{
		{
			ID:      "some_id_0",
			Title:   "some title 0",
			Formats: []ydl.Format{{FormatID: "251"}},
			Thumbnails: []ydl.Thumbnail{{
				URL:    "some url",
				ID:     "1",
				Height: 1080,
				Width:  1920,
			}},
			Description: "some description 0",
			Uploader:    "some uploader 0",
			UploaderID:  "some_uploader 0",
			UploaderURL: "some_uploader_url_0",
			UploadDate:  ydl.NewJSONTime(time.Date(2022, time.February, 27, 0, 0, 0, 0, time.UTC)),
			Filename:    "some_filename_0",
			Duration:    123,
		},
		{
			ID:      "some_id_1",
			Title:   "some title 1",
			Formats: []ydl.Format{{FormatID: "251"}},
			Thumbnails: []ydl.Thumbnail{{
				URL:    "some url",
				ID:     "1",
				Height: 1080,
				Width:  1920,
			}},
			Description: "some description 1",
			Uploader:    "some uploader 1",
			UploaderID:  "some_uploader 1",
			UploaderURL: "some_uploader_url_1",
			UploadDate:  ydl.NewJSONTime(time.Date(2022, time.February, 28, 0, 0, 0, 0, time.UTC)),
			Filename:    "some_filename_1",
			Duration:    124,
		},
	}
	results := []ydl.DownloadResult{
		{}, {},
	}
	s.Equal(len(metadata), len(results))

	order := make([]*gomock.Call, 0, len(metadata)+1)
	order = append(order, s.mockYdl.EXPECT().PlaylistMetadata(url).
		DoAndReturn(func(url string) ([]ydl.VideoMetadata, error) {
			return metadata, nil
		}))
	for i := 0; i < len(metadata); i += 1 {
		order = append(order, s.mockYdl.EXPECT().Download(metadata[i]).
			DoAndReturn(func(metadata ydl.VideoMetadata) (ydl.DownloadResult, error) {
				return results[0], nil
			}))
	}
	gomock.InOrder(order...)

	app := app.NewApplicationImpl(s.mockYdl)
	s.Nil(app.DownloadPlaylist(url))
}

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}
