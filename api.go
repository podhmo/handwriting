package handwriting

import (
	"fmt"
	"go/types"
	"io"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/podhmo/handwriting/name"
	"github.com/podhmo/handwriting/output"
	"golang.org/x/tools/go/loader"
)

// Handwriting :
type Handwriting struct {
	Pkg    *types.Package
	Config *loader.Config

	Resolver *name.Resolver
	Files    map[string]*File

	// options
	TypeCheck bool
}

// NewFromPath :
func NewFromPath(path string) *Handwriting {
	elems := strings.Split(path, "/")
	return New(types.NewPackage(path, elems[len(elems)-1]))
}

// New :
func New(pkg *types.Package, ops ...func(*Handwriting)) *Handwriting {
	h := &Handwriting{
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
	return h
}

// WithConfig :
func WithConfig(c *loader.Config) func(*Handwriting) {
	return func(h *Handwriting) {
		h.Config = c
	}
}

// importt :
func (h *Handwriting) importt(pkg string) {
	h.Config.Import(pkg)
}

// Commit :
func (h *Handwriting) Commit(w io.Writer) error {
	prog, err := h.Config.Load()
	if err != nil {
		return errors.Wrap(err, "commit")
	}
	s := &State{
		Output: output.New(w),
		Prog:   prog,
		Pkg:    h.Pkg,
	}

	files := make([]*File, 0, len(h.Files))
	for k := range h.Files {
		files = append(files, h.Files[k])
	}
	sort.Slice(files, func(i, j int) bool { return files[i].filename < files[j].filename })

	for i := range files {
		s.File = files[i]
		for _, ac := range s.File.Setups {
			if err := ac(s); err != nil {
				return errors.Wrap(err, fmt.Sprintf("setup %d, in %q", i, s.File.Name))
			}
		}
		for _, ac := range s.File.Actions {
			if err := ac(s); err != nil {
				return errors.Wrap(err, fmt.Sprintf("action %d, in %q", i, s.File.Name))
			}
		}
	}
	return nil
}

// File :
func (h *Handwriting) File(name string) *File {
	f, ok := h.Files[name]
	if !ok {
		f = &File{filename: name, Root: h, File: h.Resolver.File(nil)}
		h.Files[name] = f
	}
	return f
}

// Code :
func (f *File) Code(fn func(*State) error) {
	f.Actions = append(f.Actions, fn)
}

// Import :
func (f *File) Import(path string) {
	f.ImportWithName(path, "")
}

// ImportWithName :
func (f *File) ImportWithName(path string, name string) {
	f.Setups = append(f.Setups, func(s *State) error {
		info := s.Prog.Package(path)
		if info == nil {
			return errors.Errorf("package not found %q", path)
		}

		name := name
		if name == "" {
			name = info.Pkg.Name()
		}
		f.File.ImportWithName(info.Pkg, name)
		return nil
	})
	f.Root.importt(path)
}

// File :
type File struct {
	filename string
	*name.File
	Root    *Handwriting
	Setups  []func(*State) error
	Actions []func(*State) error
}

// State :
type State struct {
	Prog   *loader.Program
	Output *output.Output
	Pkg    *types.Package
	File   *File
}

// Lookup :
func (s *State) Lookup(name string) *types.Package {
	info := s.Prog.Package(name)
	if info == nil {
		return nil
	}
	return info.Pkg
}

// todo : sync writer
