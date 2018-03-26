package deriving

import (
	"fmt"
	"go/types"

	"github.com/pkg/errors"
	"github.com/podhmo/handwriting"
	"github.com/podhmo/handwriting/codegen/typeutil"
	"github.com/podhmo/handwriting/indent"
	"golang.org/x/tools/go/loader"
)

// BindStringer :
func BindStringer(f *handwriting.File, name string) func(e *handwriting.Emitter) error {
	f.Import("fmt")
	f.Code(func(e *handwriting.Emitter) error {
		return Stringer(e.Pkg, name, e.Output)
	})
	return nil
}

// Stringer :
func Stringer(info *loader.PackageInfo, name string, o *indent.Output) error {
	target := info.Pkg.Scope().Lookup(name)
	if target == nil {
		return errors.Errorf("%q is not found from package %q", name, info.Pkg.Path())
	}

	// todo : reuse
	typeMap := typeutil.Scan(info.Pkg)

	o.Println("// String :")
	o.WithBlock(fmt.Sprintf("func (x %s) String() string", target.Name()), func() {
		o.Println("switch x {")
		switch types.Identical(target.Type().Underlying(), types.Universe.Lookup("string").Type()) {
		case true:
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
		default:
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
		}
		o.Println("}")
	})
	return nil
}
