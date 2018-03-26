package typesutil

import (
	"go/types"
	"log"
)

// PackageDetector :
type PackageDetector struct {
	visited  []types.Type
	callback func(pkg *types.Package)
}

// NewPackageDetector :
func NewPackageDetector(callback func(pkg *types.Package)) *PackageDetector {
	return &PackageDetector{
		visited:  []types.Type{},
		callback: callback,
	}
}

// DetectPackage :
func DetectPackage(typ types.Type, callback func(pkg *types.Package)) {
	NewPackageDetector(callback).Detect(typ)
}

// Detect :
func (d *PackageDetector) Detect(typ types.Type) {
	for _, t := range d.visited {
		if t == typ {
			return
		}
	}
	d.visited = append(d.visited, typ)

	switch t := typ.(type) {
	case nil:
		return
	case *types.Named:
		d.callback(t.Obj().Pkg())
	case *types.Tuple:
		for i := 0; i < t.Len(); i++ {
			v := t.At(i)
			d.Detect(v.Type())
		}
		return
	case *types.Signature:
		{
			t := t.Params()
			for i := 0; i < t.Len(); i++ {
				v := t.At(i)
				d.Detect(v.Type())
			}

		}
		{
			t := t.Results()
			for i := 0; i < t.Len(); i++ {
				v := t.At(i)
				d.Detect(v.Type())
			}

		}
		// Params,Result
	case *types.Struct:
		if t.NumFields() > 0 {
			log.Printf("sorry raw struct expression(*types.Struct), is not supported yet in %q", d.visited[0])
		}
		// iterate field
		return
	case *types.Interface:
		if t.NumMethods() > 0 {
			log.Printf("sorry raw interface expression(*types.Interface), is not supported yet in %q", d.visited[0])
		}
		// iterate all-methods
	case *types.Map:
		d.Detect(t.Key())
		d.Detect(t.Elem())
	case hasElem:
		// *Array,*Slice,*Pointer,*Chan
		d.Detect(t.Elem())
	default: // nil, *types.Basic
		return
	}
}

type hasElem interface {
	Elem() types.Type
}
