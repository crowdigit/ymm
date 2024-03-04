package internal

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

// read reads from reader writing the read bytes to buffer.
// This function is intended to be invoked as goroutine. Any error occurs
// during read will be sent chErr and read will return. nil value means read
// operation has succeeded.
func read(name string, reader io.ReadCloser, buffer *bytes.Buffer, chErr chan<- error) {
	// bytes.Buffer.Write may panic if becomes too large
	defer func() {
		if err := recover(); err != nil {
			chErr <- fmt.Errorf("read goroutine for %s recovered: %v", name, err)
		}
	}()
	for {
		subBuffer := make([]byte, 1024)
		read, err := reader.Read(subBuffer)
		if read > 0 {
			buffer.Write(subBuffer[:read])
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				chErr <- nil
			} else {
				chErr <- err
			}
			return
		}
	}
}
