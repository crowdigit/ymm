package ydl

import (
	jsoniter "github.com/json-iterator/go"

	"time"
)

type Thumbnail struct {
	URL    string `json:"url"`
	ID     string `json:"id"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type Format struct {
	FormatID string `json:"format_id"`
}

type JSONTime time.Time

func (f *JSONTime) UnmarshalJSON(input []byte) error {
	var inputRawStr string
	if err := jsoniter.Unmarshal(input, &inputRawStr); err != nil {
		return err
	}

	t, err := time.Parse("20060102", inputRawStr)
	if err != nil {
		return err
	}

	*f = JSONTime(t)

	return nil
}

type VideoMetadata struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Formats     []Format    `json:"formats"`
	Thumbnails  []Thumbnail `json:"thumbnails"`
	Description string      `json:"description"`
	Uploader    string      `json:"uploader"`
	UploaderID  string      `json:"uploader_id"`
	UploaderURL string      `json:"uploader_url"`
	UploadDate  *JSONTime   `json:"upload_date"`
	Filename    string      `json:"filename"`
	Duration    int         `json:"duration"`
}

type DownloadResult struct{}

//go:generate mockgen -destination=../mock/mock_ydl.go -package=mock github.com/crowdigit/ymm/ydl YoutubeDL
type YoutubeDL interface {
	PlaylistMetadata(url string) ([]VideoMetadata, error)
	VideoMetadata(url string) (VideoMetadata, error)
	Download(url VideoMetadata) (DownloadResult, error)
}

type YoutubeDLImpl struct{}

func NewYoutubeDLImpl() YoutubeDL {
	return YoutubeDLImpl{}
}

func (ydl YoutubeDLImpl) PlaylistMetadata(url string) ([]VideoMetadata, error) {
	panic("not implemented") // TODO: Implement
}

func (ydl YoutubeDLImpl) VideoMetadata(url string) (VideoMetadata, error) {
	panic("not implemented") // TODO: Implement
}

func (ydl YoutubeDLImpl) Download(url VideoMetadata) (DownloadResult, error) {
	panic("not implemented") // TODO: Implement
}
