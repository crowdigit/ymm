package ydl

type VideoMetadata struct{}

//ge:generate mockgen -destination=../../mock/mock_ydl.go -package=mock github.com/crowdigit/ymm/ydl YoutubeDL
type YoutubeDL interface {
	PlaylistMetadata(url string) []VideoMetadata
	VideoMetadata(url string) VideoMetadata
	Download(url VideoMetadata)
}
