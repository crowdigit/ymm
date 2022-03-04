package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		Input  string
		Expect ProgressMessage
	}{
		{
			"[youtube] Ss-ba-g82-0: Downloading webpage\n",
			ProgressExtractor{id: "Ss-ba-g82-0"},
		},
		{
			"[download] Destination: gaburyu - エンパスィー _ 可不-Ss-ba-g82-0.webm\n",
			ProgressDestination{destination: "gaburyu - エンパスィー _ 可不-Ss-ba-g82-0.webm"},
		},
		{
			"[download]   1.0% of 2.91MiB at 120.99KiB/s ETA 00:24",
			ProgressDownload{percentage: 10},
		},
		{
			"[download] 100% of 2.91MiB in 00:38\n",
			ProgressDownload{percentage: 1000},
		},
		{
			"[ffmpeg] Destination: gaburyu - エンパスィー _ 可不-Ss-ba-g82-0.mp3\n",
			ProgressFFMPEGDestination{destination: "gaburyu - エンパスィー _ 可不-Ss-ba-g82-0.mp3"},
		},
		{
			"Deleting original file gaburyu - エンパスィー _ 可不-Ss-ba-g82-0.webm (pass -k to keep)\n",
			ProgressFFMPEGDelete{},
		},
	}

	for _, test := range tests {
		result, err := parseLine(test.Input)
		assert.Nil(t, err)
		assert.Equal(t, result, test.Expect)
	}
}

func TestParseStream(t *testing.T) {
	testInput, err := os.Open("./test.txt")
	assert.Nil(t, err)
	defer testInput.Close()

	chMessage, err := parseStream(testInput)
	assert.Nil(t, err)

	progresses := make([]ProgressMessage, 0, 49)
	for progress := range chMessage {
		progresses = append(progresses, progress)
	}

	percentages := []int{0, 1, 2, 5, 10, 21, 43, 66, 90, 118, 140, 170, 198, 225, 249, 275, 300, 328, 350, 377, 403, 428, 459, 484, 507, 536, 564, 587, 614, 642, 669, 697, 721, 748, 769, 796, 821, 848, 872, 901, 929, 952, 979, 1000, 1000}
	expected := make([]ProgressMessage, 0, 49)
	expected = append(expected, ProgressExtractor{})
	expected = append(expected, ProgressDestination{})
	for _, percentage := range percentages {
		expected = append(expected, ProgressDownload{
			percentage: percentage,
		})
	}
	expected = append(expected, ProgressFFMPEGDestination{})
	expected = append(expected, ProgressFFMPEGDelete{})

	assert.Equal(t, len(expected), len(progresses))
	assert.Equal(t, expected, progresses)
}
