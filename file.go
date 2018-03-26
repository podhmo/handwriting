package handwriting

import (
	"github.com/pkg/errors"
	"github.com/podhmo/handwriting/name"
)

// File :
type File struct {
	Filename string
	*name.File
	Root *Planner

	Headers []string
	Setups  []func(*Emitter) error
	Actions []func(*Emitter) error
	imports []importspec
	used    map[string]struct{}
}

// Header :
func (f *File) Header(s string) {
	f.Headers = append(f.Headers, s)
}

// Code :
func (f *File) Code(fn func(*Emitter) error) {
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

	skipimport := false
	for _, pkgspec := range f.Root.Config.CreatePkgs {
		if pkgspec.Path == path {
			skipimport = true
			break
		}
	}
	if !skipimport {
		f.Root.Config.Import(path)
	}

	if f.Root.Pkg.Path() == path {
		return
	}

	f.Setups = append(f.Setups, func(s *Emitter) error {
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
	f.imports = append(f.imports, importspec{Name: name, Path: path})

}

type importspec struct {
	Name string
	Path string
}
