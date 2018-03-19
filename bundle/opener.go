package bundle

import (
	"io"
)

// Opener :
type Opener interface {
	Open(name string) (io.WriteCloser, error)
}
