package deriving

import (
	"fmt"
	"go/types"

	"github.com/pkg/errors"
	"github.com/podhmo/handwriting"
	"github.com/podhmo/handwriting/codegen/lookup"
	"github.com/podhmo/handwriting/codegen/typesutil"
	"github.com/podhmo/handwriting/indent"
)

// BindStringer :
func BindStringer(f *handwriting.File, name string) func(e *handwriting.Emitter) error {
	f.Import("fmt")
	f.Code(func(e *handwriting.Emitter) error {
		return Stringer(e.PkgInfo.Pkg, name, e.Output)
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
		for _, ob := range typeMap[target.Type()] {
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
		for _, ob := range typeMap[target.Type()] {
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
