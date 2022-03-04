package main

import (
	"fmt"
	"io"
)

type Progress struct {
	Percentage int
}

type ProgressMessage interface {
	Update(*Progress)
}

type ProgressExtractor struct{}

func (p ProgressExtractor) Update(*Progress) {
}

type ProgressDestination struct{}

func (p ProgressDestination) Update(*Progress) {
}

type ProgressDownload struct {
	percentage int
}

func (p ProgressDownload) Update(*Progress) {
}

type ProgressFFMPEGDestination struct{}

func (p ProgressFFMPEGDestination) Update(*Progress) {
}

type ProgressFFMPEGDelete struct{}

func (p ProgressFFMPEGDelete) Update(*Progress) {
}

func parseStream(read io.Reader) (chan ProgressMessage, error) {
	chMessage := make(chan ProgressMessage)
	defer close(chMessage)
	return chMessage, nil
}

func main() {
	fmt.Println("vim-go")
}
