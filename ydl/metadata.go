package ydl

import (
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
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

func NewJSONTime(input time.Time) *JSONTime {
	output := JSONTime(input)
	return &output
}

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

func (f *JSONTime) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("%q", (*time.Time)(f).Format("20060102"))
	return []byte(str), nil
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
	WebpageURL  string      `json:"webpage_url"`
	Filename    string      `json:"_filename"`
	Duration    int         `json:"duration"`
}
