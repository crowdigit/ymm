package main

import (
	"log"

	"github.com/crowdigit/ymm/ydl"
)

func main() {
	commandProvider := ydl.NewCommandProviderImpl()
	youtubeDl := ydl.NewYoutubeDLImpl(commandProvider)

	url := "https://www.youtube.com/watch?v=Ss-ba-g82-0"
	metadata, err := youtubeDl.VideoMetadata(url)
	if err != nil {
		log.Fatalf("failed to get video metadata: %s", err)
	}

	log.Printf("%v", metadata)
}
