package ydl

import "io"

//go:generate mockgen -destination=../mock/mock_command.go -package=mock github.com/crowdigit/ymm/ydl CommandProvider
type CommandProvider interface {
	Start(name string, arg ...string) error
	StderrPipe() (io.ReadCloser, error)
	StdoutPipe() (io.ReadCloser, error)
	Wait() (int, error)
}
