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
	Download(metadata VideoMetadata) (DownloadResult, error)
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

func readStream(wg *sync.WaitGroup, reader io.ReadCloser, chOut chan<- []byte, chErr chan<- error) {
	defer wg.Done()
	readBuffer := make([]byte, 1024)
	for {
		read, err := reader.Read(readBuffer)
		if read > 0 {
			sendBuffer := make([]byte, read)
			copy(sendBuffer, readBuffer)
			chOut <- sendBuffer
		} else if err == io.EOF {
			break
		} else if err != nil {
			chErr <- errors.Wrap(err, "failed to read from reader stream")
			break
		}
	}
}

func handleMetadataStream(chStdout <-chan []byte, chStderr <-chan []byte, chJson chan<- []byte, chClose chan struct{}, chErr <-chan error) {
	json := make([]byte, 0, 8192)
loop:
	for {
		select {
		case <-chErr:
			// TODO
			break loop
		case <-chStderr:
			// TODO
		case jsonChunk := <-chStdout:
			json = append(json, jsonChunk...)
		case <-chClose:
			break loop
		}
	}
	chJson <- json
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

	chStderr := make(chan []byte)
	chStdout := make(chan []byte)
	chErr := make(chan error)
	go readStream(&wg, stderr, chStderr, chErr)
	go readStream(&wg, stdout, chStdout, chErr)

	if err := command.Start(); err != nil {
		return VideoMetadata{}, errors.Wrap(err, "failed to start metadata command")
	}

	chClose := make(chan struct{})
	chJson := make(chan []byte)
	go handleMetadataStream(chStdout, chJson, chClose)

	wg.Wait()

	close(chClose)
	json := <-chJson

	status, err := command.Wait()
	if err != nil {
		return VideoMetadata{}, errors.Wrap(err, "failed to wait for metadata command")
	}

	if status != 0 {
		return VideoMetadata{}, fmt.Errorf("metadata command exited with %d", status)
	}

	result := VideoMetadata{}
	if err := jsoniter.Unmarshal(json, &result); err != nil {
		return VideoMetadata{}, errors.Wrap(err, "failed to unmarshal bytes from stdout")
	}

	return result, nil
}

func (ydl YoutubeDLImpl) Download(metadata VideoMetadata) (DownloadResult, error) {
	panic("not implemented") // TODO: Implement
}
