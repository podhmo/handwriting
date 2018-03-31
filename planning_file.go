package handwriting

import (
	"github.com/pkg/errors"
)

// PlanningFile :
type PlanningFile struct {
	Filename string
	Root     *Planner

	Setups  []func(*File) error
	Actions []func(*File) error
	imports []importspec
	used    map[string]struct{}
}

// Code :
func (f *PlanningFile) Code(fn func(*File) error) {
	f.Actions = append(f.Actions, fn)
}

// Import :
func (f *PlanningFile) Import(path string) {
	f.ImportWithName(path, "")
}

// ImportWithName :
func (f *PlanningFile) ImportWithName(path string, name string) {
	if _, existed := f.used[path]; existed {
		return
	}
	f.used[path] = struct{}{}
	f.Root.Import(path)
	if f.Root.Pkg.Path() == path {
		return
	}

	f.Setups = append(f.Setups, func(f *File) error {
		info := f.Prog.Package(path)
		if info == nil {
			return errors.Errorf("package not found %q", path)
		}

		name := name
		if name == "" {
			name = info.Pkg.Name()
		}
		f.Resolver.ImportWithName(info.Pkg, name)
		return nil
	})
	f.imports = append(f.imports, importspec{Name: name, Path: path})

}

type importspec struct {
	Name string
	Path string
}
