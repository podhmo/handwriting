package transform

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/pkg/errors"
	"github.com/podhmo/handwriting"
	"github.com/podhmo/handwriting/codegen/lookup"
	"github.com/podhmo/handwriting/codegen/typesutil"
	"github.com/podhmo/handwriting/indent"
)

// EmitAsFakeStruct :
func EmitAsFakeStruct(f *handwriting.File, path string, exportedOnly bool) func(e *handwriting.Emitter) error {
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

	// todo : comment
	outname := fmt.Sprintf("Fake%s", name)
	o.Printfln("// %s is fake struct of %s", outname, types.TypeString(target.Type(), types.RelativeTo(f.Root.Pkg)))

	n := iface.NumMethods()
	o.WithBlock(fmt.Sprintf("type %s struct", outname), func() {
		// import pkg, if not imported yet.
		d := typesutil.NewPackageDetector(func(pkg *types.Package) {
			if pkg != nil {
				f.Import(pkg.Path())
			}
		})

		for i := 0; i < n; i++ {
			method := iface.Method(i)
			if exportedOnly && !method.Exported() {
				continue
			}
			d.Detect(method.Type())
			o.Printf("%s%s %s", strings.ToLower(method.Name()[0:1]), method.Name()[1:], f.TypeName(method.Type()))
		}
	})

	for i := 0; i < n; i++ {
		method := iface.Method(i)
		if exportedOnly && !method.Exported() {
			continue
		}
		sig, _ := method.Type().(*types.Signature)
		if sig == nil {
			return errors.New("hmm") // xxx :
		}
		o.WithBlock(fmt.Sprintf("func %s %s %s", method.Name(), sig.Params(), sig.Results()), func() {
		})
	}
	return nil
}
