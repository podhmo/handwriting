package handwriting

import (
	"go/types"
	"strings"

	"github.com/podhmo/handwriting/multifile"
	"golang.org/x/tools/go/loader"
)

// NewFromPackagePath :
func NewFromPackagePath(path string, ops ...func(*Planner)) (*Planner, error) {
	elems := strings.Split(path, "/")
	pkg := types.NewPackage(path, elems[len(elems)-1])
	return New(pkg, ops...)
}

// WithConfig :
func WithConfig(c *loader.Config) func(*Planner) {
	return func(h *Planner) {
		h.Config = c
	}
}

// WithOpener :
func WithOpener(o multifile.Opener) func(*Planner) {
	return func(h *Planner) {
		h.Opener = o
	}
}

// WithDryRun :
func WithDryRun() func(*Planner) {
	return WithOpener(multifile.Stderr())
}
