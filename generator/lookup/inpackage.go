package lookup

import (
	"fmt"
	"go/types"
	"log"

	"github.com/podhmo/handwriting/generator/typesutil"
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
func (ref *StructRef) IterateMethods(fn func(*types.Func), mode typesutil.IterateMode) {
	named := ref.Named
	n := named.NumMethods()
	for i := 0; i < n; i++ {
		method := named.Method(i)
		if mode == typesutil.All || (method.Exported() && mode == typesutil.ExportedOnly) {
			fn(method)
		}
	}
}

// IterateAllMethods :
func (ref *StructRef) IterateAllMethods(fn func(*types.Func)) {
	ref.IterateMethods(fn, typesutil.All)
}

// IterateFields :
func (ref *StructRef) IterateFields(fn func(*types.Var), mode typesutil.IterateMode) {
	s := ref.Underlying
	n := s.NumFields()
	for i := 0; i < n; i++ {
		field := s.Field(i)
		if mode == typesutil.All || (field.Exported() && mode == typesutil.ExportedOnly) {
			fn(field)
		}
	}
}

// IterateAllFields :
func (ref *StructRef) IterateAllFields(fn func(*types.Var)) {
	ref.IterateFields(fn, typesutil.All)
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

// IterateMethods :
func (ref *InterfaceRef) IterateMethods(fn func(*FuncRef), mode typesutil.IterateMode) {
	iface := ref.Underlying
	n := iface.NumMethods()
	for i := 0; i < n; i++ {
		method := iface.Method(i)
		if mode == typesutil.All || (method.Exported() && mode == typesutil.ExportedOnly) {
			sig, _ := method.Type().(*types.Signature)
			if sig == nil {
				log.Printf("invalid method? %q doensn't have signature", method.Name())
				continue
			}
			fn(&FuncRef{Func: method, Signature: sig})
		}
	}
}

// IterateAllMethods :
func (ref *InterfaceRef) IterateAllMethods(fn func(*FuncRef)) {
	ref.IterateMethods(fn, typesutil.All)
}

// FuncRef :
type FuncRef struct {
	*types.Func
	Signature *types.Signature
}

// AsFunc :
func AsFunc(ob types.Object) (*FuncRef, error) {
	fn, _ := ob.(*types.Func)
	if fn == nil {
		return nil, &lookupError{Type: Type("func"), Msg: fmt.Sprintf("%q is not func", ob.Name()), Where: ob.Pkg().Path()}
	}
	signature, _ := fn.Type().(*types.Signature)
	if signature == nil {
		return &FuncRef{Func: fn, Signature: signature}, nil
	}
	return &FuncRef{Func: fn, Signature: signature}, nil
}
