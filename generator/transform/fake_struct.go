package transform

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/pkg/errors"
	"github.com/podhmo/handwriting"
	"github.com/podhmo/handwriting/generator/lookup"
	"github.com/podhmo/handwriting/generator/typesutil"
	"github.com/podhmo/handwriting/nameresolve"
)

// TODO : struct's name policy
// TODO : the subject on method definition, name policy

// GenerateFakeStruct :
func GenerateFakeStruct(f *handwriting.PlanningFile, path string, exportedOnly bool) error {
	// path = <package path>/<name>
	elems := strings.Split(path, "/")
	pkgpath := strings.Join(elems[:len(elems)-1], "/")
	name := elems[len(elems)-1]

	f.Import(pkgpath)
	f.Code(func(f *handwriting.File) error {
		g := GeneratorForFakeStructNew(f)
		pkg, err := g.f.Use(pkgpath)
		if err != nil {
			return errors.Wrap(err, "lookup package")
		}
		iface, err := pkg.LookupInterface(name)
		if err != nil {
			return errors.Wrap(err, "lookup interface")
		}
		return g.Generate(iface, name, "fake"+name, exportedOnly)
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

// GeneratorForFakeStruct :
type GeneratorForFakeStruct struct {
	f *handwriting.File
	d *typesutil.PackageDetector
}

// GeneratorForFakeStructNew :
func GeneratorForFakeStructNew(f *handwriting.File) *GeneratorForFakeStruct {
	return &GeneratorForFakeStruct{
		f: f,
		d: f.CreateCaptureImportDetector(),
	}
}

// Generate :
func (g *GeneratorForFakeStruct) Generate(iface *lookup.InterfaceRef, name, outname string, exportedOnly bool) error {
	// todo : comment
	o := g.f.Out
	r := g.f.Resolver
	o.Printfln("// %s is fake struct of %s", outname, name)

	// define struct
	o.WithBlock(fmt.Sprintf("type %s struct", outname), func() {
		iface.IterateMethods(typesutil.IterateModeFromBool(exportedOnly), func(method *lookup.FuncRef) {
			g.d.Detect(method.Type())
			o.Printfln("%s %s", nameresolve.ToUnexported(method.Name()), r.TypeName(method.Type()))
		})
	})

	// define methods
	iface.IterateMethods(typesutil.IterateModeFromBool(exportedOnly), func(method *lookup.FuncRef) {
		sig := method.Signature
		o.Printfln("// %s :", method.Name())

		// xxx : fill params name
		params := sig.Params()
		unnamed := false
		for i := 0; i < params.Len(); i++ {
			x := params.At(i)
			if x.Name() == "" {
				unnamed = true
				break
			}
		}
		if unnamed {
			namedVars := make([]*types.Var, params.Len())
			for i := range namedVars {
				x := params.At(i)
				if x.Name() != "" {
					namedVars[i] = x
					continue
				}
				namedVars[i] = types.NewVar(x.Pos(), x.Pkg(), fmt.Sprintf("v%d", i), x.Type())
			}
			params = types.NewTuple(namedVars...)
		}

		o.WithBlock(fmt.Sprintf("func (x *%s) %s %s %s", outname, method.Name(), r.TypeName(params), r.TypeNameForResults(sig.Results())), func() {
			params := params

			varnames := make([]string, params.Len())
			for i := 0; i < params.Len(); i++ {
				x := params.At(i)
				varnames[i] = x.Name()
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
