package multifile

import (
	"go/types"
	"io"
	"os"
)

// Opener :
type Opener interface {
	Open(name string) (io.WriteCloser, error)
}

// WriteFile :
func WriteFile(opener Opener, filename string, fn func(w io.Writer) error) error {
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

// Stdout :
func Stdout() Opener {
	return &consoleOpener{w: os.Stdout}
}

// Stderr :
func Stderr() Opener {
	return &consoleOpener{w: os.Stderr}
}

// Dir :
func Dir(base string) (Opener, error) {
	if err := os.MkdirAll(base, 0744); err != nil {
		return nil, err
	}
	return &fileOpener{Base: base}, nil
}

// Package :
func Package(pkg *types.Package, createIfNotExists bool) (Opener, error) {
	path, err := pkgFilePath(pkg.Path(), createIfNotExists)
	if err != nil {
		return nil, err
	}
	return &fileOpener{Base: path}, nil
}

// Must :
func Must(o Opener, err error) Opener {
	if err != nil {
		panic(err)
	}
	return o
}
