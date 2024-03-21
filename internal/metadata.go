package internal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/crowdigit/exec"
	jsoniter "github.com/json-iterator/go"
)

type JqMetadata struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Formats []struct {
		FormatID string `json:"format_id"`
	} `json:"formats"`
	Thumbnail   string   `json:"thumbnail"`
	Description string   `json:"description"`
	ChannelID   string   `json:"channel_id"`
	ChannelURL  string   `json:"channel_url"`
	Duration    int      `json:"duration"`
	WebpageURL  string   `json:"webpage_url"`
	Categories  []string `json:"categories"`
	Tags        []string `json:"tags"`
	Chapters    []struct {
		StartTime string `json:"start_time"`
		Title     string `json:"title"`
		EndTime   string `json:"end_time"`
	} `json:"chapters"`
	Channel            string `json:"channel"`
	Uploader           string `json:"uploader"`
	UploaderID         string `json:"uploader_id"`
	UploaderURL        string `json:"uploader_url"`
	UploadDate         string `json:"upload_date"`
	OriginalURL        string `json:"original_url"`
	WebpageURLBasename string `json:"webpage_url_basename"`
	WebpageURLDomain   string `json:"webpage_url_domain"`
	Extractor          string `json:"extractor"`
	DisplayID          string `json:"display_id"`
	Fulltitle          string `json:"fulltitle"`
	Epoch              int64  `json:"epoch"`
	Filename           string `json:"filename"`
}

func fetchMetadata(
	cp exec.CommandProvider,
	youtube, json ConfigExec,
	url string,
) ([]JqMetadata, error) {
	ctx, kill := context.WithCancel(context.Background())
	defer kill()
	ctx, unregister := signal.NotifyContext(ctx, os.Interrupt)
	defer unregister()

	youtube = youtube.ReplacePlaceholder("<url>", url)

	buffers := []*bytes.Buffer{{}, {}}
	pipespecs := []exec.PipeSpec{
		{
			CmdOpt: exec.CommandOpts{
				Path: youtube.Path,
				Args: youtube.Args,
			},
			Next:  exec.Stdout,
			Other: buffers[0],
		},
		{
			CmdOpt: exec.CommandOpts{Path: json.Path, Args: json.Args},
			Next:   exec.Stdout,
			Other:  buffers[1],
		},
	}

	pipeline, err := exec.NewPipeline(ctx, cp, pipespecs)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize pipeline: %w", err)
	}

	err = pipeline.Start()
	defer pipeline.Cancel(kill)
	if err != nil {
		return nil, fmt.Errorf("failed to start pipeline: %w", err)
	}

	var metadata []JqMetadata
	chPipeErr := make(chan error)
	go func() {
		defer close(chPipeErr)
		decoder := jsoniter.NewDecoder(pipeline.Output())
		if err := decoder.Decode(&metadata); err != nil && !errors.Is(err, io.EOF) {
			chPipeErr <- err
			return
		}
	}()
	if err := <-chPipeErr; err != nil {
		return nil, fmt.Errorf("failed to operate on pipe: %w", err)
	}

	if errs := pipeline.Any(kill); errs != nil {
		fmt.Fprintln(os.Stderr, errs)
		for i, err := range errs {
			fmt.Println(err)
			if buffers[errs[i].Index].Len() > 0 {
				fmt.Println(buffers[i].String())
			}
		}
		return nil, errs
	}

	return metadata, nil
}

func selectFormat(metadata JqMetadata) (string, error) {
	switch metadata.Extractor {
	case "youtube":
		// Youtube - 251 format is required
		for _, format := range metadata.Formats {
			if format.FormatID == "251" {
				return format.FormatID, nil
			}
		}
		// TODO error with context
		return "", fmt.Errorf("required format for Youtube video (251) is missing")
	default:
		return "", fmt.Errorf("not implemented extractor: %s", metadata.Extractor)
	}
}
