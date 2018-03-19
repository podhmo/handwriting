package handwriting

import (
	"go/types"
	"os"
	"strings"

	"github.com/podhmo/handwriting/opener"
	"golang.org/x/tools/go/loader"
)

// NewFromPackagePath :
func NewFromPackagePath(path string, ops ...func(*Handwriting)) (*Handwriting, error) {
	elems := strings.Split(path, "/")
	pkg := types.NewPackage(path, elems[len(elems)-1])
	return New(pkg, ops...)
}

// WithConfig :
func WithConfig(c *loader.Config) func(*Handwriting) {
	return func(h *Handwriting) {
		h.Config = c
	}
}

// WithOpener :
func WithOpener(o opener.Opener) func(*Handwriting) {
	return func(h *Handwriting) {
		h.Opener = o
	}
}

// WithDryRun :
func WithDryRun() func(*Handwriting) {
	return WithOpener(opener.Console(os.Stderr))
}
