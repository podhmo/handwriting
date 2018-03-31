package lookup

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/loader"
)

// Package :
func Package(prog *loader.Program, pkgpath string) (*PackageRef, error) {
	info := prog.Package(pkgpath)
	if info == nil {
		return nil, &lookupError{Type: Type("package"), Msg: fmt.Sprintf("%q is not found", pkgpath)}
	}
	return &PackageRef{info.Pkg}, nil
}

// PackageRef :
type PackageRef struct {
	*types.Package
}

// LookupStruct :
func (ref *PackageRef) LookupStruct(name string) (*StructRef, error) {
	ob, err := Object(ref.Package, name)
	if err != nil {
		return nil, err
	}
	return AsStruct(ob)
}

// LookupInterface :
func (ref *PackageRef) LookupInterface(name string) (*InterfaceRef, error) {
	ob, err := Object(ref.Package, name)
	if err != nil {
		return nil, err
	}
	return AsInterface(ob)
}

// LookupFunc :
func (ref *PackageRef) LookupFunc(name string) (*FuncRef, error) {
	ob, err := Object(ref.Package, name)
	if err != nil {
		return nil, err
	}
	return AsFunc(ob)
}
