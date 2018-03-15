package handwriting

import (
	"go/ast"
	"go/types"
)

// todo: move it

// ImportedNameTo :
func ImportedNameTo(pkg *types.Package, f *ast.File) types.Qualifier {
	return func(other *types.Package) string {
		if pkg == other {
			return "" // same package; unqualified
		}
		// todo: cache
		for _, is := range f.Imports {
			if is.Path.Value[1:len(is.Path.Value)-1] == other.Path() {
				if is.Name != nil {
					return is.Name.String()
				}
				return other.Name()
			}
		}
		return other.Name() // todo: add import
	}
}
