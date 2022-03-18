package ydl

type VideoMetadata struct{}

//go:generate mockgen -destination=../mock/mock_ydl.go -package=mock github.com/crowdigit/ymm/ydl YoutubeDL
type YoutubeDL interface {
	PlaylistMetadata(url string) ([]VideoMetadata, error)
	VideoMetadata(url string) (VideoMetadata, error)
	Download(url VideoMetadata) error
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

func (ydl YoutubeDLImpl) Download(url VideoMetadata) error {
	panic("not implemented") // TODO: Implement
}
