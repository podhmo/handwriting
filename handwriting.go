package handwriting

import (
	"fmt"
	"go/types"
	"sort"

	"github.com/pkg/errors"
	"github.com/podhmo/handwriting/name"
	"github.com/podhmo/handwriting/opener"
	"github.com/podhmo/handwriting/output"
	"golang.org/x/tools/go/loader"
)

// Handwriting :
type Handwriting struct {
	Pkg    *types.Package
	Config *loader.Config

	Resolver *name.Resolver
	Files    map[string]*File
	Opener   opener.Opener
	// options
	TypeCheck bool
}

// New :
func New(pkg *types.Package, ops ...func(*Handwriting)) (*Handwriting, error) {
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
	if h.Opener == nil {
		createIfNotExists := true
		opener, err := opener.NewFromPackage(pkg, createIfNotExists)
		if err != nil {
			return nil, err
		}
		h.Opener = opener
	}
	return h, nil
}

// importt :
func (h *Handwriting) importt(pkg string) {
	h.Config.Import(pkg)
}

// Commit :
func (h *Handwriting) Commit() error {
	prog, err := h.Config.Load()
	if err != nil {
		return errors.Wrap(err, "commit")
	}
	s := &State{
		Prog: prog,
		Pkg:  h.Pkg,
	}

	files := make([]*File, 0, len(h.Files))
	for k := range h.Files {
		files = append(files, h.Files[k])
	}
	sort.Slice(files, func(i, j int) bool { return files[i].filename < files[j].filename })

	for i := range files {
		i := i
		s.File = files[i]
		w, err := h.Opener.Open(s.File.filename)
		if err != nil {
			return err
		}
		if err := func() error {
			defer w.Close()
			s.Output = output.New(w)
			for _, ac := range s.File.Setups {
				if err := ac(s); err != nil {
					return errors.Wrap(err, fmt.Sprintf("setup %d, in %q", i, s.File.Name))
				}
			}

			s.Output.Printf("package %s\n", s.Pkg.Name())
			if len(s.File.imports) > 0 {
				s.Output.Println("")
				s.Output.Println("import (")
				s.Output.Indent()
				// todo : sort
				for _, im := range s.File.imports {
					if im.Name == "" {
						s.Output.Printfln(`%q`, im.Path)
					} else {
						s.Output.Printfln(`%s %q`, im.Name, im.Path)
					}
				}
				s.Output.UnIndent()
				s.Output.Println(")")
				s.Output.Println("")
			}

			for _, ac := range s.File.Actions {
				if err := ac(s); err != nil {
					return errors.Wrap(err, fmt.Sprintf("action %d, in %q", i, s.File.Name))
				}
			}
			return nil
		}(); err != nil {
			return err
		}
	}
	return nil
}

// File :
func (h *Handwriting) File(name string) *File {
	f, ok := h.Files[name]
	if !ok {
		f = &File{filename: name, Root: h, File: h.Resolver.File(nil), used: map[string]struct{}{}}
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
	if _, existed := f.used[path]; existed {
		return
	}
	f.used[path] = struct{}{}

	f.Setups = append(f.Setups, func(s *State) error {
		info := s.Prog.Package(path)
		if info == nil {
			return errors.Errorf("package not found %q", path)
		}

		name := name
		f.imports = append(f.imports, importspec{Name: name, Path: path})
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
	Root *Handwriting

	Setups  []func(*State) error
	Actions []func(*State) error
	imports []importspec
	used    map[string]struct{}
}

type importspec struct {
	Name string
	Path string
}

// State :
type State struct {
	Prog *loader.Program
	Pkg  *types.Package
	*output.Output
	*File
}

// Lookup :
func (s *State) Lookup(name string) *types.Package {
	info := s.Prog.Package(name)
	if info == nil {
		return nil
	}
	return info.Pkg
}
