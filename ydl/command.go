package ydl

import "io"

//go:generate mockgen -destination=../mock/mock_command.go -package=mock github.com/crowdigit/ymm/ydl Command
type Command interface {
	Start() error
	StderrPipe() (io.ReadCloser, error)
	StdoutPipe() (io.ReadCloser, error)
	Wait() (int, error)
}

//go:generate mockgen -destination=../mock/mock_command_provider.go -package=mock github.com/crowdigit/ymm/ydl CommandProvider
type CommandProvider interface {
	NewCommand(name string, arg ...string) Command
}
