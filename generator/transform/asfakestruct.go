package transform

import (
	"fmt"
	"go/types"
	"log"
	"strings"

	"github.com/pkg/errors"
	"github.com/podhmo/handwriting"
	"github.com/podhmo/handwriting/generator/lookup"
	"github.com/podhmo/handwriting/generator/typesutil"
	"github.com/podhmo/handwriting/indent"
	"github.com/podhmo/handwriting/nameresolve"
)

// TODO : struct's name policy
// TODO : the subject on method definition, name policy

// GenerateFakeStruct :
func GenerateFakeStruct(f *handwriting.File, path string, exportedOnly bool) func(e *handwriting.Emitter) error {
	// path = <package path>/<name>
	elems := strings.Split(path, "/")
	pkgpath := strings.Join(elems[:len(elems)-1], "/")
	name := elems[len(elems)-1]

	f.Import(pkgpath)
	f.Code(func(e *handwriting.Emitter) error {
		info, err := lookup.PackageInfo(e.Prog, pkgpath)
		if info == nil {
			return errors.Wrap(err, "lookup pacakge")
		}
		return AsFakeStruct(f, info.Pkg, name, e.Output, exportedOnly)
	})
	return nil
}

/*
type I interface {
	F(x string) string
}

// To

type FakeI struct {
	f func(x string) string
}

func (x *FakeI) F(x string) string {
	return x.f(x)
}
*/

// AsFakeStruct :
func AsFakeStruct(f *handwriting.File, pkg *types.Package, name string, o *indent.Output, exportedOnly bool) error {
	target, err := lookup.Object(pkg, name)
	if err != nil {
		return errors.Wrap(err, "lookup target")
	}
	iface, err := lookup.AsInterface(target)
	if err != nil {
		return errors.Wrap(err, "lookup interface")
	}

	// import pkg, if not imported yet.
	d := typesutil.NewPackageDetector(func(pkg *types.Package) {
		if pkg != nil {
			f.Import(pkg.Path())
		}
	})

	// todo : comment
	outname := fmt.Sprintf("Fake%s", name)
	o.Printfln("// %s is fake struct of %s", outname, types.TypeString(target.Type(), types.RelativeTo(f.Root.Pkg)))

	// define struct
	o.WithBlock(fmt.Sprintf("type %s struct", outname), func() {
		iface.IterateMethods(typesutil.IterateModeFromBool(exportedOnly), func(method *types.Func) {
			d.Detect(method.Type())
			o.Printfln("%s %s", nameresolve.ToUnexported(method.Name()), f.TypeName(method.Type()))
		})
	})

	// define methods
	iface.IterateMethods(typesutil.IterateModeFromBool(exportedOnly), func(method *types.Func) {
		sig, _ := method.Type().(*types.Signature)
		if sig == nil {
			log.Printf("invalid method? %q doensn't have signature", method.Name())
			return
		}
		o.Printfln("// %s :", method.Name())
		o.WithBlock(fmt.Sprintf("func (x *%s) %s %s %s", outname, method.Name(), f.TypeName(sig.Params()), f.TypeNameForResults(sig.Results())), func() {
			params := sig.Params()
			varnames := make([]string, params.Len())
			for i := 0; i < params.Len(); i++ {
				varnames[i] = params.At(i).Name()
			}

			if sig.Results().Len() == 0 {
				o.Printfln("x.%s(%s)", nameresolve.ToUnexported(method.Name()), strings.Join(varnames, ", "))
				return
			}
			o.Printfln("return x.%s(%s)", nameresolve.ToUnexported(method.Name()), strings.Join(varnames, ", "))
		})
	})
	return nil
}