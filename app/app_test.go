package app_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/crowdigit/ymm/app"
	"github.com/crowdigit/ymm/db"
	"github.com/crowdigit/ymm/mock"
	"github.com/crowdigit/ymm/ydl"
	"github.com/golang/mock/gomock"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/suite"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"go.uber.org/zap"
)

type AppTestSuite struct {
	suite.Suite

	mockCtrl *gomock.Controller
	mockYdl  *mock.MockYoutubeDL
	mockDb   *mock.MockDatabase
	bundb    *bun.DB
}

func (s *AppTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockYdl = mock.NewMockYoutubeDL(s.mockCtrl)
	s.mockDb = mock.NewMockDatabase(s.mockCtrl)

	sqldb, err := sql.Open(sqliteshim.ShimName, ":memory:")
	s.Nil(err)
	s.bundb = bun.NewDB(sqldb, sqlitedialect.New())
}

func (s *AppTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
	s.bundb.Close()
}

func (s *AppTestSuite) TestDownloadSingle() {
	url := "http://asdf.com/some/url"
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
		UploadDate:  ydl.NewJSONTime(time.Date(2022, time.February, 27, 0, 0, 0, 0, time.UTC)),
		Filename:    "some_filename",
		Duration:    123,
	}
	metadataBytes, err := jsoniter.Marshal(metadata)
	s.Nil(err)
	result := ydl.DownloadResult{}

	order := []*gomock.Call{
		s.mockYdl.EXPECT().VideoMetadata(url).
			Return(metadataBytes, nil),
		s.mockDb.EXPECT().BunDB().
			Return(s.bundb),
		s.mockDb.EXPECT().SelectUploader(gomock.Any()).
			Return([]db.Uploader{}, nil),
		s.mockDb.EXPECT().BunDB().
			Return(s.bundb),
		s.mockDb.EXPECT().InsertUploader(gomock.Any()).
			Return(nil),
		s.mockDb.EXPECT().StoreMetadata(metadata.ID, metadataBytes).
			Return(nil),
		s.mockYdl.EXPECT().Download(gomock.Any(), metadata).
			Return(result, nil),
	}
	gomock.InOrder(order...)

	app := app.NewApplicationImpl(zap.NewNop().Sugar(), s.mockYdl, s.mockDb)
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
	metadataBytes := make([][]byte, 0, len(metadata))
	for _, metadatum := range metadata {
		metadatumBytes, err := jsoniter.Marshal(metadatum)
		s.Nil(err)
		metadataBytes = append(metadataBytes, metadatumBytes)
	}
	results := []ydl.DownloadResult{
		{}, {},
	}
	s.Equal(len(metadata), len(results))

	order := make([]*gomock.Call, 0, len(metadata)*2+1)
	order = append(order, s.mockYdl.EXPECT().PlaylistMetadata(url).
		DoAndReturn(func(url string) ([][]byte, error) {
			return metadataBytes, nil
		}))
	for i, metadatumBytes := range metadataBytes {
		order = append(order, s.mockDb.EXPECT().BunDB().
			Return(s.bundb))
		order = append(order, s.mockDb.EXPECT().SelectUploader(gomock.Any()).
			Return([]db.Uploader{}, nil))
		order = append(order, s.mockDb.EXPECT().BunDB().
			Return(s.bundb))
		order = append(order, s.mockDb.EXPECT().InsertUploader(gomock.Any()).
			Return(nil))
		order = append(order, s.mockDb.EXPECT().StoreMetadata(metadata[i].ID, metadatumBytes).
			Return(nil))
	}
	for i := 0; i < len(metadata); i += 1 {
		order = append(order, s.mockYdl.EXPECT().Download(gomock.Any(), metadata[i]).
			Return(results[i], nil))
	}
	gomock.InOrder(order...)

	app := app.NewApplicationImpl(zap.NewNop().Sugar(), s.mockYdl, s.mockDb)
	s.Nil(app.DownloadPlaylist(url))
}

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}
