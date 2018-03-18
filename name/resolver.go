package name

import (
	"fmt"
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
	named := map[string]int{}
	if f != nil {
		for _, is := range f.Imports {
			path := is.Path.Value[1 : len(is.Path.Value)-1]
			if is.Name != nil {
				name := is.Name.String()
				imported[path] = name
				named[name] = 1
			}
		}
	}
	return &File{Root: r, Imported: imported, named: named}
}

// File :
type File struct {
	Root     *Resolver
	Imported map[string]string // path -> prefix
	named    map[string]int
	fakes    map[string]*types.Package
}

// Import :
func (f *File) Import(pkg *types.Package) {
	f.ImportWithName(pkg, pkg.Name())
	return
}

// ImportWithName :
func (f *File) ImportWithName(pkg *types.Package, name string) {
	if _, ok := f.Imported[pkg.Path()]; ok {
		return
	}
	i, ok := f.named[name]
	if !ok {
		f.Imported[pkg.Path()] = name
		f.named[name]++
		return
	}
	f.Imported[pkg.Path()] = fmt.Sprintf("%s%d", name, i)
	f.named[name]++
	return
}

// Name :
func (f *File) Name(ob types.Object) string {
	prefix := f.Prefix(ob.Pkg())
	if prefix == "" {
		return ob.Name()
	}
	return fmt.Sprintf("%s.%s", prefix, ob.Name())
}

// TypeName :
func (f *File) TypeName(typ types.Type) string {
	return types.TypeString(typ, f.Prefix)
}

// Prefix :
func (f *File) Prefix(other *types.Package) string {
	if f.Root.Pkg == other {
		return "" // same package; unqualified
	}
	path := other.Path()
	if name, ok := f.Imported[path]; ok {
		return name
	}
	f.Import(other)
	return f.Imported[other.Path()]
}

// shorthand

// ImportFake :
func ImportFake(f *File, path, name string) {
	if _, ok := f.Imported[path]; ok {
		return
	}
	pkg := types.NewPackage(path, name)
	f.fakes[path] = pkg
	f.ImportWithName(pkg, name)
	return
}

// WithPrefix :
func WithPrefix(f *File, path, name string) string {
	ImportFake(f, path, name)
	prefix := f.Prefix(f.fakes[path])
	if prefix == "" {
		return name
	}
	return fmt.Sprintf("%s.%s", prefix, name)
}