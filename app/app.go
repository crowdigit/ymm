package app

import "github.com/crowdigit/ymm/ydl"

//go:generate mockgen -destination=../mock/mock_app.go -package=mock github.com/crowdigit/ymm/app Application
type Application interface {
	DownloadPlaylist(url string)
	DownloadSingle(url string)
}

type ApplicationImpl struct {
	ydl ydl.YoutubeDL
}

func NewApplicationImpl(ydl ydl.YoutubeDL) Application {
	return ApplicationImpl{}
}

func (app ApplicationImpl) DownloadPlaylist(url string) {
	panic("not implemented") // TODO: Implement
}

func (app ApplicationImpl) DownloadSingle(url string) {
	panic("not implemented") // TODO: Implement
}
