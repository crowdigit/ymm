package ydl

type VideoMetadata struct{}

//ge:generate mockgen -destination=../../mock/mock_ydl.go -package=mock github.com/crowdigit/ymm/ydl YoutubeDL
type YoutubeDL interface {
	PlaylistMetadata(url string) []VideoMetadata
	VideoMetadata(url string) VideoMetadata
	Download(url VideoMetadata)
}

type YoutubeDLImpl struct{}

func NewYoutubeDLImpl() YoutubeDL {
	return YoutubeDLImpl{}
}

func (ydl YoutubeDLImpl) PlaylistMetadata(url string) []VideoMetadata {
	panic("not implemented") // TODO: Implement
}

func (ydl YoutubeDLImpl) VideoMetadata(url string) VideoMetadata {
	panic("not implemented") // TODO: Implement
}

func (ydl YoutubeDLImpl) Download(url VideoMetadata) {
	panic("not implemented") // TODO: Implement
}
