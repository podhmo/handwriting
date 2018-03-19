package bundle

import (
	"go/build"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// NewFromPackage :
func NewFromPackage(pkg *types.Package, createIfNotExists bool) (Opener, error) {
	path, err := pkgFilePath(pkg.Path(), createIfNotExists)
	if err != nil {
		return nil, err
	}
	return &fileOpener{Base: path}, nil
}

func pkgFilePath(pkgname string, createIfNotExists bool) (string, error) {
	ctxt := build.Default
	for _, srcdir := range ctxt.SrcDirs() {
		path := filepath.Join(srcdir, pkgname)
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			return path, nil
		}
	}
	if !createIfNotExists {
		return "", errors.Errorf("%q's physical path is not found", pkgname)
	}
	for _, srcdir := range ctxt.SrcDirs() {
		if strings.HasPrefix(srcdir, ctxt.GOPATH) {
			path := filepath.Join(srcdir, pkgname)
			return path, os.MkdirAll(path, 0744)
		}
	}
	return "", errors.Errorf("%q's physical path is not found", pkgname)
}
