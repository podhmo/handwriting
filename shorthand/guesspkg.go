package shorthand

import (
	"go/build"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// GuessPkg :
func GuessPkg() (string, error) {
	curdir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return GuessPkgFromPath(curdir)
}

// GuessPkgFromPath :
func GuessPkgFromPath(curdir string) (string, error) {
	path, err := filepath.Abs(curdir)
	if err != nil {
		return "", err
	}
	for _, srcdir := range build.Default.SrcDirs() {
		if strings.HasPrefix(path, srcdir) {
			pkgname := strings.TrimLeft(strings.Replace(path, srcdir, "", 1), "/")
			return pkgname, nil
		}
	}
	return "", errors.Errorf("%q is not subdir of srcdirs(%q)", path, build.Default.SrcDirs())
}
