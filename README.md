# YMM

## Features

1. Download music from Youtube
   * Download a video file as audio
   * Download a playlist as audio files
2. Generate ReplayGain (EBU R128) MP3 tag on each file
3. Manage downloaded files
   * Make structured directory by creator
   * Put files in its creators' directory
   * Let user change creator directory's name (not implemented yet)

## Examples

1. Download a video as mp3 `example/single/main.go`
2. Download a playlist `example/playlist/main.go`

Downloaded files will be located at `$XDG_DATA_HOME/ymm/music`

## How to download private video/playlist

1. Make a merged cookie file
   * [See](https://daveparrish.net/posts/2018-06-22-How-to-download-private-YouTube-videos-with-youtube-dl.html)
2. Pass path to cookie file to `NewYoutubeDLImpl`
