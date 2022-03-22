package ydl

import (
	"fmt"
	"io"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

type DownloadResult struct{}

//go:generate mockgen -destination=../mock/mock_ydl.go -package=mock github.com/crowdigit/ymm/ydl YoutubeDL
type YoutubeDL interface {
	PlaylistMetadata(url string) ([]VideoMetadata, error)
	VideoMetadata(url string) (VideoMetadata, error)
	Download(url VideoMetadata) (DownloadResult, error)
}

type YoutubeDLImpl struct {
	commandProvider CommandProvider
}

func NewYoutubeDLImpl(command CommandProvider) YoutubeDL {
	return YoutubeDLImpl{
		commandProvider: command,
	}
}

func (ydl YoutubeDLImpl) PlaylistMetadata(url string) ([]VideoMetadata, error) {
	panic("not implemented") // TODO: Implement
}

func (ydl YoutubeDLImpl) VideoMetadata(url string) (VideoMetadata, error) {
	command := ydl.commandProvider.NewCommand("youtube-dl", "--dump-json", url)

	stderr, err := command.StderrPipe()
	if err != nil {
		return VideoMetadata{}, errors.Wrap(err, "failed to get stderr pipe")
	}

	stdout, err := command.StdoutPipe()
	if err != nil {
		return VideoMetadata{}, errors.Wrap(err, "failed to get stdout pipe")
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		buffer := make([]byte, 1024)
		for {
			read, err := stderr.Read(buffer)
			if read > 0 {
				// TODO log error
			} else if err == io.EOF {
				break
			} else if err != nil {
				// TODO log error
				// TODO return error
				break
			}
		}
		wg.Done()
	}()

	chChunk := make(chan []byte)

	go func() {
		buffer := make([]byte, 1024)
		for {
			read, err := stdout.Read(buffer)
			if read > 0 {
				sendbuffer := make([]byte, read)
				copy(sendbuffer, buffer)
				chChunk <- sendbuffer
			} else if err == io.EOF {
				break
			} else if err != nil {
				// TODO log error
				// TODO return error
				break
			}
		}
		close(chChunk)
		wg.Done()
	}()

	if err := command.Start(); err != nil {
		return VideoMetadata{}, errors.Wrap(err, "failed to start metadata command")
	}

	json := make([]byte, 0, 8192)
	for chunk := range chChunk {
		json = append(json, chunk...)
	}

	status, err := command.Wait()
	if err != nil {
		return VideoMetadata{}, errors.Wrap(err, "failed to wait for metadata command")
	}

	if status != 0 {
		return VideoMetadata{}, fmt.Errorf("metadata command exited with %d", status)
	}

	wg.Wait()

	result := VideoMetadata{}
	if err := jsoniter.Unmarshal(json, &result); err != nil {
		return VideoMetadata{}, errors.Wrap(err, "failed to unmarshal bytes from stdout")
	}

	return result, nil
}

func (ydl YoutubeDLImpl) Download(url VideoMetadata) (DownloadResult, error) {
	panic("not implemented") // TODO: Implement
}
