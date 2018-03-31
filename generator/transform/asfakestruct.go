package transform

import (
	"fmt"
	"go/types"
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
func GenerateFakeStruct(f *handwriting.PlanningFile, path string, exportedOnly bool) func(e *handwriting.Emitter) error {
	// path = <package path>/<name>
	elems := strings.Split(path, "/")
	pkgpath := strings.Join(elems[:len(elems)-1], "/")
	name := elems[len(elems)-1]

	f.Import(pkgpath)
	f.Code(func(f *handwriting.File) error {
		pkg, err := lookup.Package(f.Prog, pkgpath)
		if pkg == nil {
			return errors.Wrap(err, "lookup pacakge")
		}
		return AsFakeStruct(f, pkg, name, f.Out, exportedOnly)
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
	d := f.CreateCaptureImportDetector()

	// todo : comment
	outname := fmt.Sprintf("Fake%s", name)
	o.Printfln("// %s is fake struct of %s", outname, types.TypeString(target.Type(), types.RelativeTo(f.PkgInfo.Pkg)))

	// define struct
	o.WithBlock(fmt.Sprintf("type %s struct", outname), func() {
		iface.IterateMethods(typesutil.IterateModeFromBool(exportedOnly), func(method *lookup.FuncRef) {
			d.Detect(method.Type())
			o.Printfln("%s %s", nameresolve.ToUnexported(method.Name()), f.Resolver.TypeName(method.Type()))
		})
	})

	// define methods
	iface.IterateMethods(typesutil.IterateModeFromBool(exportedOnly), func(method *lookup.FuncRef) {
		sig := method.Signature
		o.Printfln("// %s :", method.Name())
		o.WithBlock(fmt.Sprintf("func (x *%s) %s %s %s", outname, method.Name(), f.Resolver.TypeName(sig.Params()), f.Resolver.TypeNameForResults(sig.Results())), func() {
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
