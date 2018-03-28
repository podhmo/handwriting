package lookup

import (
	"fmt"
	"go/types"
)

// Object :
func Object(pkg *types.Package, name string) (types.Object, error) {
	ob := pkg.Scope().Lookup(name)
	if ob == nil {
		return nil, &lookupError{Type: Type("object"), Msg: fmt.Sprintf("%q is not found", name), Where: pkg.Path()}
	}
	return ob, nil
}

// StructRef :
type StructRef struct {
	*types.Named
	Underlying *types.Struct
}

// AsStruct :
func AsStruct(ob types.Object) (*StructRef, error) {
	named, _ := ob.Type().(*types.Named)
	if named == nil {
		return nil, &lookupError{Type: Type("named"), Msg: fmt.Sprintf("%q is not struct", ob.Name()), Where: ob.Pkg().Path()}
	}
	underlying, _ := named.Underlying().(*types.Struct)
	if underlying == nil {
		return nil, &lookupError{Type: Type("struct"), Msg: fmt.Sprintf("%q is not struct", ob.Name()), Where: ob.Pkg().Path()}
	}
	return &StructRef{Named: named, Underlying: underlying}, nil
}

// IterateMethods :
func (ref *StructRef) IterateMethods(fn func(*types.Func)) {
	n := ref.NumMethods()
	for i := 0; i < n; i++ {
		method := ref.Method(i)
		fn(method)
	}
}

// InterfaceRef :
type InterfaceRef struct {
	*types.Named
	Underlying *types.Interface
}

// AsInterface :
func AsInterface(ob types.Object) (*InterfaceRef, error) {
	named, _ := ob.Type().(*types.Named)
	if named == nil {
		return nil, &lookupError{Type: Type("ob.Name()d"), Msg: fmt.Sprintf("%q is not ob.Name()d", ob.Name()), Where: ob.Pkg().Path()}
	}
	underlying, _ := named.Underlying().(*types.Interface)
	if underlying == nil {
		return nil, &lookupError{Type: Type("interface"), Msg: fmt.Sprintf("%q is not interface", ob.Name()), Where: ob.Pkg().Path()}
	}
	return &InterfaceRef{Named: named, Underlying: underlying}, nil
}
