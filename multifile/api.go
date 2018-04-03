package multifile

import (
	"go/build"
	"go/types"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
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
		return nil, errors.Wrap(err, "create default directory")
	}
	if err := cleaning(base); err != nil {
		return nil, errors.Wrap(err, "cleaning files")
	}
	return &fileOpener{Base: base}, nil
}

func cleaning(dirpath string) error {
	fs, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return err
	}

	for _, f := range fs {
		if f.Size() == 0 {
			log.Printf("cleaning, remove %s (empty file)", f.Name())
			if err := os.Remove(filepath.Join(dirpath, f.Name())); err != nil {
				return err
			}
		}
	}
	return nil
}

// Package :
func Package(pkg *types.Package, createIfNotExists bool) (Opener, error) {
	if build.IsLocalImport(pkg.Path()) {
		return Dir(pkg.Path())
	}

	path, err := pkgFilePath(pkg.Path(), createIfNotExists)
	if err != nil {
		return nil, err
	}
	return Dir(path)
}

// Must :
func Must(o Opener, err error) Opener {
	if err != nil {
		panic(err)
	}
	return o
}
