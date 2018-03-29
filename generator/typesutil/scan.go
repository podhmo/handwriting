package typesutil

import (
	"go/types"
	"reflect"
)

// Scan :
func Scan(pkg *types.Package) map[types.Type][]types.Object {
	r := map[types.Type][]types.Object{}
	s := pkg.Scope()
	for _, name := range reflect.ValueOf(s).Elem().FieldByName("elems").MapKeys() {
		ob := s.Lookup(name.String())
		r[ob.Type()] = append(r[ob.Type()], ob)
	}
	return r
}
