package handwriting

import (
	"go/types"
	"sort"

	"github.com/pkg/errors"
	"github.com/podhmo/handwriting/multifile"
	"github.com/podhmo/handwriting/name"
	"golang.org/x/tools/go/loader"
)

// Planner :
type Planner struct {
	Pkg    *types.Package
	Config *loader.Config

	Resolver *name.Resolver
	Files    map[string]*File
	Opener   multifile.Opener
	// options
	TypeCheck bool
}

// New :
func New(pkg *types.Package, ops ...func(*Planner)) (*Planner, error) {
	h := &Planner{
		Pkg:      pkg,
		Resolver: name.New(pkg),
		Files:    map[string]*File{},
	}

	for _, op := range ops {
		op(h)
	}

	if h.Config == nil {
		h.Config = &loader.Config{}
		if !h.TypeCheck {
			h.Config.TypeCheckFuncBodies = func(path string) bool {
				return false
			}
		}
	}
	if h.Opener == nil {
		createIfNotExists := true
		opener, err := multifile.Package(pkg, createIfNotExists)
		if err != nil {
			return nil, err
		}
		h.Opener = opener
	}
	return h, nil
}

// Emit :
func (h *Planner) Emit() error {
	prog, err := h.Config.Load()
	if err != nil {
		return errors.Wrap(err, "commit")
	}
	r := &Emitter{
		Prog:   prog,
		Pkg:    h.Pkg,
		Opener: h.Opener,
	}

	files := make([]*File, 0, len(h.Files))
	for k := range h.Files {
		files = append(files, h.Files[k])
	}
	sort.Slice(files, func(i, j int) bool { return files[i].filename < files[j].filename })

	for i := range files {
		if err := r.Emit(files[i]); err != nil {
			return err
		}
	}
	return nil
}

// File :
func (h *Planner) File(name string) *File {
	f, ok := h.Files[name]
	if !ok {
		f = &File{filename: name, Root: h, File: h.Resolver.File(nil), used: map[string]struct{}{}}
		h.Files[name] = f
	}
	return f
}
