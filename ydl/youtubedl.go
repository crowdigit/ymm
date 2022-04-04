package ydl

import (
	"fmt"
	"io"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type DownloadResult struct{}

//go:generate mockgen -destination=../mock/mock_ydl.go -package=mock github.com/crowdigit/ymm/ydl YoutubeDL
type YoutubeDL interface {
	PlaylistMetadata(url string) ([][]byte, error)
	VideoMetadata(url string) ([]byte, error)
	Download(workDir string, metadata VideoMetadata) (DownloadResult, error)
}

type YoutubeDLImpl struct {
	logger          *zap.SugaredLogger
	commandProvider CommandProvider
}

func NewYoutubeDLImpl(logger *zap.SugaredLogger, command CommandProvider) YoutubeDL {
	return YoutubeDLImpl{
		logger:          logger,
		commandProvider: command,
	}
}

func contains251Format(formats []Format) bool {
	for _, format := range formats {
		if format.FormatID == "251" {
			return true
		}
	}
	return false
}

func (ydl YoutubeDLImpl) PlaylistMetadata(url string) ([][]byte, error) {
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

func (ydl YoutubeDLImpl) VideoMetadata(url string) ([]byte, error) {
	command := ydl.commandProvider.NewCommand("youtube-dl", "--dump-json", url)

	stderr, err := command.StderrPipe()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stderr pipe")
	}

	stdout, err := command.StdoutPipe()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stdout pipe")
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	chStderr := make(chan []byte)
	chStdout := make(chan []byte)
	chErr := make(chan error)
	go readStream(&wg, stderr, chStderr, chErr)
	go readStream(&wg, stdout, chStdout, chErr)

	if err := command.Start(); err != nil {
		return nil, errors.Wrap(err, "failed to start metadata command")
	}

	chClose := make(chan struct{})
	chJson := make(chan []byte)
	go handleMetadataStream(chStdout, chStderr, chJson, chClose, chErr)

	wg.Wait()

	close(chClose)
	json := <-chJson

	status, err := command.Wait()
	if err != nil {
		return nil, errors.Wrap(err, "failed to wait for metadata command")
	}

	if status != 0 {
		return nil, fmt.Errorf("metadata command exited with %d", status)
	}

	return json, nil
}

func handleDownloadStream(logger *zap.SugaredLogger, chStdout <-chan []byte, chStderr <-chan []byte, chClose chan struct{}, chErr <-chan error) {
loop:
	for {
		select {
		case err := <-chErr:
			logger.Errorf("reading from download stream returned an error: %s", err)
			break loop
		case msg := <-chStderr:
			logger.Warnf(string(msg))
		case msg := <-chStdout:
			logger.Infof(string(msg))
		case <-chClose:
			break loop
		}
	}
}

func (ydl YoutubeDLImpl) Download(workDir string, metadata VideoMetadata) (DownloadResult, error) {
	if !contains251Format(metadata.Formats) {
		return DownloadResult{}, fmt.Errorf("video %s does not contain format 251", metadata.ID)
	}

	command := ydl.commandProvider.NewCommand(
		"youtube-dl",
		"--format", "251",
		"--extract-audio",
		"--audio-format", "mp3",
		"--audio-quality", "0",
		metadata.WebpageURL)

	stderr, err := command.StderrPipe()
	if err != nil {
		return DownloadResult{}, errors.Wrap(err, "failed to get stderr pipe")
	}

	stdout, err := command.StdoutPipe()
	if err != nil {
		return DownloadResult{}, errors.Wrap(err, "failed to get stdout pipe")
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	chStderr := make(chan []byte)
	chStdout := make(chan []byte)
	chErr := make(chan error)
	go readStream(&wg, stderr, chStderr, chErr)
	go readStream(&wg, stdout, chStdout, chErr)

	if err := command.Start(); err != nil {
		return DownloadResult{}, errors.Wrap(err, "failed to start download command")
	}

	chClose := make(chan struct{})
	go handleDownloadStream(ydl.logger, chStdout, chStderr, chClose, chErr)

	wg.Wait()

	close(chClose)

	status, err := command.Wait()
	if err != nil {
		return DownloadResult{}, errors.Wrap(err, "failed to wait for download command")
	}

	if status != 0 {
		return DownloadResult{}, fmt.Errorf("download command exited with %d", status)
	}

	return DownloadResult{}, nil
}
