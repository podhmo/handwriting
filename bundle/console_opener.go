package bundle

import (
	"fmt"
	"io"
)

// Console :
func Console(w io.Writer) Opener {
	return &consoleOpener{w: w}
}

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
