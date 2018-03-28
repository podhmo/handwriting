package typesutil

import "go/types"

// IterateMode :
type IterateMode int

// IterateMode :
const (
	ExportedOnly IterateMode = iota
	All
)

// IterateObjects :
func IterateObjects(pkg *types.Package, mode IterateMode, fn func(types.Object)) {
	s := pkg.Scope()
	for _, name := range s.Names() {
		ob := s.Lookup(name)
		if mode == All || (ob.Exported() && mode == ExportedOnly) {
			fn(ob)
		}
	}
}

// IterateAllObjects :
func IterateAllObjects(pkg *types.Package, fn func(types.Object)) {
	IterateObjects(pkg, All, fn)
}

// IterateModeFromBool :
func IterateModeFromBool(exportedOnly bool) IterateMode {
	if exportedOnly {
		return ExportedOnly
	}
	return All
}
