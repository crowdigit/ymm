package ydl_test

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/crowdigit/ymm/ydl"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

var commonTestMetadata = ydl.VideoMetadata{
	ID:    "Ss-ba-g82-0",
	Title: "gaburyu - エンパスィー /可不",
	Formats: []ydl.Format{
		{FormatID: "250"},
		{FormatID: "251"},
	},
	Thumbnails: []ydl.Thumbnail{
		{
			URL:    "https://i.ytimg.com/vi/Ss-ba-g82-0/hqdefault.jpg?sqp=-oaymwEcCNACELwBSFXyq4qpAw4IARUAAIhCGAFwAcABBg==&rs=AOn4CLDu_ZbPvpCCP5x3RpSxFLJRweAyzg",
			ID:     "3",
			Height: 188,
			Width:  336,
		},
		{
			URL:    "https://i.ytimg.com/vi_webp/Ss-ba-g82-0/maxresdefault.webp",
			ID:     "4",
			Height: 1080,
			Width:  1920,
		},
	},
	Description: "qwer",
	Uploader:    "gaburyu",
	UploaderID:  "UC-V8o6aa9C133RYMtaobq6w",
	UploaderURL: "http://www.youtube.com/channel/UC-V8o6aa9C133RYMtaobq6w",
	UploadDate:  ydl.NewJSONTime(time.Date(2022, time.February, 27, 0, 0, 0, 0, time.UTC)),
	WebpageURL:  "https://www.youtube.com/watch?v=Ss-ba-g82-0",
	Filename:    "gaburyu - エンパスィー _可不-Ss-ba-g82-0.mp4",
	Duration:    174,
}

func TestMetadataUnmarshal(t *testing.T) {
	inputFile, err := os.Open("./test.json")
	assert.Nil(t, err)
	defer inputFile.Close()

	input, err := io.ReadAll(inputFile)
	assert.Nil(t, err)

	got := ydl.VideoMetadata{}
	assert.Nil(t, jsoniter.Unmarshal(input, &got))

	expected := commonTestMetadata
	assert.Equal(t, expected, got)
}

var commonTestMetadataList = []ydl.VideoMetadata{
	{
		ID:      "E6GimcGTCvk",
		Title:   "半透明 / 初音ミク MV",
		Formats: []ydl.Format{{FormatID: "251"}},
		Thumbnails: []ydl.Thumbnail{
			{
				URL:    "https://i.ytimg.com/vi_webp/E6GimcGTCvk/maxresdefault.webp",
				ID:     "4",
				Height: 1080,
				Width:  1920,
			},
		},
		Description: "asdf",
		Uploader:    "獅子志司",
		UploaderID:  "UCihQ9WPpwc-fZ9lz3ZjAYnA",
		UploaderURL: "http://www.youtube.com/channel/UCihQ9WPpwc-fZ9lz3ZjAYnA",
		UploadDate:  ydl.NewJSONTime(time.Date(2021, time.July, 31, 0, 0, 0, 0, time.UTC)),
		WebpageURL:  "https://www.youtube.com/watch?v=E6GimcGTCvk",
		Filename:    "半透明 _ 初音ミク MV-E6GimcGTCvk.webm",
		Duration:    197,
	},
	{
		ID:      "228WfzixKVw",
		Title:   "鬣犬新書 / 初音ミク MV",
		Formats: []ydl.Format{{FormatID: "251"}},
		Thumbnails: []ydl.Thumbnail{{
			URL:    "https://i.ytimg.com/vi/228WfzixKVw/hqdefault.jpg?sqp=-oaymwEcCNACELwBSFXyq4qpAw4IARUAAIhCGAFwAcABBg==&rs=AOn4CLBDF4taNUQmyl11-NKSWch7Tr86OA",
			ID:     "3",
			Height: 188,
			Width:  336,
		}},
		Description: "zxcv",
		Uploader:    "獅子志司",
		UploaderID:  "UCihQ9WPpwc-fZ9lz3ZjAYnA",
		UploaderURL: "http://www.youtube.com/channel/UCihQ9WPpwc-fZ9lz3ZjAYnA",
		UploadDate:  &ydl.JSONTime{},
		WebpageURL:  "https://www.youtube.com/watch?v=228WfzixKVw",
		Filename:    "鬣犬新書 _ 初音ミク MV-228WfzixKVw.mp4",
		Duration:    223,
	},
}
