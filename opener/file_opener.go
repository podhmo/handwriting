package opener

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

type fileOpener struct {
	Base string
}

// Open :
func (o *fileOpener) Open(filename string) (io.WriteCloser, error) {
	path := filepath.Join(o.Base, filename)
	log.Printf("open %q", path)
	return os.Create(path)
}
