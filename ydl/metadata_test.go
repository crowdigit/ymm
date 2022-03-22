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

func TestMetadataUnmarshal(t *testing.T) {
	inputFile, err := os.Open("./test.json")
	assert.Nil(t, err)
	defer inputFile.Close()

	input, err := io.ReadAll(inputFile)
	assert.Nil(t, err)

	got := ydl.VideoMetadata{}
	assert.Nil(t, jsoniter.Unmarshal(input, &got))

	expectedUploadeDate := ydl.JSONTime(time.Date(2022, time.February, 27, 0, 0, 0, 0, time.UTC))
	expected := ydl.VideoMetadata{
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
		Description: "03\n惑わされてしまう！\n\n作曲/CG : gaburyu\nVocal : 可不\n-------------------------------------\n▼カシ・オトハ ver.\nhttps://youtu.be/DcKQz72s0OM\n\n▼ほかの曲/Other Works\nhttps://youtube.com/playlist?list=PLmnvJgQL8hd5hyXYVoP9WqUGHyqZDRZ9N\n\n▼gaburyu\n❏Twitter \nhttps://twitter.com/gabustep\n❏Youtube\nhttps://www.youtube.com/c/gaburyu\n❏niconico\nhttps://www.nicovideo.jp/series/99575\n-------------------------------------\n\n\n #gaburyu #可不 #エンパスィー",
		Uploader:    "gaburyu",
		UploaderID:  "UC-V8o6aa9C133RYMtaobq6w",
		UploaderURL: "http://www.youtube.com/channel/UC-V8o6aa9C133RYMtaobq6w",
		UploadDate:  &expectedUploadeDate,
		Filename:    "gaburyu - エンパスィー _可不-Ss-ba-g82-0.mp4",
		Duration:    174,
	}
	assert.Equal(t, expected, got)
}
