package main

import (
	"fmt"
	"io"
)

const BUFFER_SIZE = 128

type Progress struct {
	ID               string
	Percentage       int
	VideoDestination string
	AudioDestination string
}

type ProgressMessage interface {
	Update(*Progress)
}

type ProgressExtractor struct {
	id string
}

func (m ProgressExtractor) Update(p *Progress) {
	p.ID = m.id
}

type ProgressDestination struct {
	destination string
}

func (m ProgressDestination) Update(p *Progress) {
	p.VideoDestination = m.destination
}

type ProgressDownload struct {
	percentage int
}

func (m ProgressDownload) Update(p *Progress) {
	p.Percentage = m.percentage
}

type ProgressFFMPEGDestination struct {
	destination string
}

func (m ProgressFFMPEGDestination) Update(p *Progress) {
	p.AudioDestination = m.destination
}

type ProgressFFMPEGDelete struct{}

func (p ProgressFFMPEGDelete) Update(*Progress) {
}

func parseLine(input string) (ProgressMessage, error) {
	return nil, nil
}

func parseStream(read io.Reader) (chan ProgressMessage, error) {
	chMessage := make(chan ProgressMessage)
	go func() {
		defer close(chMessage)
		buffer := make([]byte, BUFFER_SIZE)
		for {
			read, err := read.Read(buffer)
			if read > 0 {
				// TODO parse
			}
			if err == io.EOF {
				return
			} else if err != nil {
				return
			}
		}
	}()
	return chMessage, nil
}

func main() {
	fmt.Println("vim-go")
}
