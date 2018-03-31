package namesutil

import (
	"fmt"
	"go/types"
)

// ToNamedTuple :
func ToNamedTuple(tuple *types.Tuple) *types.Tuple {
	named := true
	for i := 0; i < tuple.Len(); i++ {
		x := tuple.At(i)
		if x.Name() == "" {
			named = false
			break
		}
	}
	if named {
		return tuple
	}

	namedVars := make([]*types.Var, tuple.Len())
	for i := range namedVars {
		x := tuple.At(i)
		if x.Name() != "" {
			namedVars[i] = x
			continue
		}
		namedVars[i] = types.NewVar(x.Pos(), x.Pkg(), fmt.Sprintf("v%d", i), x.Type())
	}
	return types.NewTuple(namedVars...)
}

// NamesFromTuple :
func NamesFromTuple(tuple *types.Tuple) []string {
	names := make([]string, tuple.Len())
	for i := 0; i < tuple.Len(); i++ {
		x := tuple.At(i)
		names[i] = x.Name()
	}
	return names
}
