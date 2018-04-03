package deriving

import (
	"fmt"
	"go/types"
	"sort"

	"github.com/pkg/errors"
	"github.com/podhmo/handwriting"
	"github.com/podhmo/handwriting/generator/lookup"
	"github.com/podhmo/handwriting/generator/typesutil"
	"github.com/podhmo/handwriting/indent"
)

// GenerateStringer :
func GenerateStringer(f *handwriting.PlanningFile, name string) error {
	f.Import("fmt")
	f.Code(func(f *handwriting.File) error {
		return Stringer(f.PkgInfo.Pkg, name, f.Out)
	})
	return nil
}

// Stringer :
func Stringer(pkg *types.Package, name string, o *indent.Output) error {
	target, err := lookup.Object(pkg, name)
	if err != nil {
		return errors.Wrap(err, "lookup target")
	}

	if types.Identical(target.Type().Underlying(), types.Universe.Lookup("string").Type()) {
		return stringerForStringType(target, pkg, name, o)
	}

	return stringerDefault(target, pkg, name, o)
}

func stringerForStringType(target types.Object, pkg *types.Package, name string, o *indent.Output) error {
	// todo : reuse
	typeMap := typesutil.Scan(pkg)
	o.Println("// String :")
	o.WithBlock(fmt.Sprintf("func (x %s) String() string", target.Name()), func() {
		o.Println("switch x {")
		candidates := typeMap[target.Type()]
		sort.Slice(candidates, func(i, j int) bool { return candidates[i].Id() < candidates[j].Id() })
		for _, ob := range candidates {
			if ob, ok := ob.(*types.Const); ok {
				o.WithIndent(fmt.Sprintf("case %s:", ob.Name()), func() {
					o.Println(fmt.Sprintf("return %s", ob.Val()))
				})
			}
		}
		o.WithIndent("default:", func() {
			o.Printfln(`return fmt.Sprintf("%s(%%q)", string(x))`, target.Name())
		})
		o.Println("}")
	})
	return nil
}

func stringerDefault(target types.Object, pkg *types.Package, name string, o *indent.Output) error {
	// todo : reuse
	typeMap := typesutil.Scan(pkg)

	o.Println("// String :")
	o.WithBlock(fmt.Sprintf("func (x %s) String() string", target.Name()), func() {
		o.Println("switch x {")
		candidates := typeMap[target.Type()]
		sort.Slice(candidates, func(i, j int) bool { return candidates[i].Id() < candidates[j].Id() })
		for _, ob := range candidates {
			if ob, ok := ob.(*types.Const); ok {
				o.WithIndent(fmt.Sprintf("case %s:", ob.Name()), func() {
					o.Println(fmt.Sprintf("return %q", ob.Name()))
				})
			}
		}
		o.WithIndent("default:", func() {
			o.Printfln(`return fmt.Sprintf("%s(%%q)", x)`, target.Name())
		})
		o.Println("}")
	})
	return nil
}
