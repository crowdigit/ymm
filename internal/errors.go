package internal

import "fmt"

type PipelineError struct {
	err   error
	Index int
}

func (e PipelineError) Error() string {
	return fmt.Sprintf("pipeline process %d has failed: %s", e.Index, e.err)
}

func (e PipelineError) Unwrap() error {
	return e.err
}
