package bundle

import (
	"go/types"
	"io"
)

// Opener :
type Opener interface {
	Open(name string) (io.WriteCloser, error)
}

// Bundle :
func Bundle(opener Opener, filename string, fn func(w io.Writer) error) error {
	w, err := opener.Open(filename)
	if err != nil {
		return err
	}
	defer w.Close()
	return fn(w)
}

// Console :
func Console(w io.Writer) Opener {
	return &consoleOpener{w: w}
}

// Package :
func Package(pkg *types.Package, createIfNotExists bool) (Opener, error) {
	path, err := pkgFilePath(pkg.Path(), createIfNotExists)
	if err != nil {
		return nil, err
	}
	return &fileOpener{Base: path}, nil
}
