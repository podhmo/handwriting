package handwriting

import (
	"github.com/pkg/errors"
	"github.com/podhmo/handwriting/name"
)

// File :
type File struct {
	filename string
	*name.File
	Root *Planner

	Setups  []func(*Emitter) error
	Actions []func(*Emitter) error
	imports []importspec
	used    map[string]struct{}
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

	f.Setups = append(f.Setups, func(s *Emitter) error {
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

type importspec struct {
	Name string
	Path string
}
