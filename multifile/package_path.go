package multifile

import (
	"go/build"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

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
