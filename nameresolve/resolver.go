package nameresolve

import (
	"go/ast"
	"go/types"
)

// Resolver :
type Resolver struct {
	Pkg *types.Package
}

// New :
func New(pkg *types.Package) *Resolver {
	return &Resolver{
		Pkg: pkg,
	}
}

// File :
func (r *Resolver) File(f *ast.File) *File {
	imported := map[string]string{}
	used := map[string]int{}
	if f != nil {
		for _, is := range f.Imports {
			path := is.Path.Value[1 : len(is.Path.Value)-1]
			if is.Name != nil {
				name := is.Name.String()
				imported[path] = name
				used[name] = 1
			}
		}
	}
	return &File{Root: r, Imported: imported, used: used}
}
