package name

import (
	"fmt"
	"go/types"
)

// File :
type File struct {
	Root     *Resolver
	Imported map[string]string // path -> prefix
	used     map[string]int
	fakes    map[string]*types.Package
}

// Import :
func (f *File) Import(pkg *types.Package) {
	f.ImportWithName(pkg, pkg.Name())
}

// ImportWithName :
func (f *File) ImportWithName(pkg *types.Package, name string) {
	if _, ok := f.Imported[pkg.Path()]; ok {
		return
	}
	importedName := name
	if i, ok := f.used[name]; !ok {
		importedName = fmt.Sprintf("%s%d", name, i)
	}
	f.Imported[pkg.Path()] = importedName
	f.used[name]++
}

// ImportFake :
func (f *File) ImportFake(path, name string) {
	if _, ok := f.Imported[path]; ok {
		return
	}
	pkg := types.NewPackage(path, name)
	f.fakes[path] = pkg
	f.ImportWithName(pkg, name)
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

// WithPrefix :
func WithPrefix(f *File, path, name string) string {
	f.ImportFake(path, name)
	prefix := f.Prefix(f.fakes[path])
	if prefix == "" {
		return name
	}
	return fmt.Sprintf("%s.%s", prefix, name)
}
