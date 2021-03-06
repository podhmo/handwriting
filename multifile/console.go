package multifile

import (
	"fmt"
	"io"
)

type consoleOpener struct {
	w io.Writer
}

// Open :
func (o *consoleOpener) Open(filename string) (io.WriteCloser, error) {
	fmt.Fprintln(o.w, "open ", filename)
	fmt.Fprintln(o.w, "----------------------------------------")
	return &nopCloser{Writer: o.w}, nil
}

// ioutil.NopCloser 's writer version
type nopCloser struct {
	io.Writer
}

// Close :
func (c *nopCloser) Close() error {
	return nil
}
