package app

//ge:generate mockgen -destination=../../mock/mock_app.go -package=mock github.com/crowdigit/ymm/app Application
type Application interface {
	DownloadPlaylist(url string)
	DownloadSingle(url string)
}
