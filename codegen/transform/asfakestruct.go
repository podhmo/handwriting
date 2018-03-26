package transform

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/pkg/errors"
	"github.com/podhmo/handwriting"
	"github.com/podhmo/handwriting/codegen/typesutil"
	"github.com/podhmo/handwriting/indent"
	"golang.org/x/tools/go/loader"
)

// EmitAsFakeStruct :
func EmitAsFakeStruct(f *handwriting.File, path string, exportedOnly bool) func(e *handwriting.Emitter) error {
	// <package path>/<name>
	elems := strings.Split(path, "/")
	pkgpath := strings.Join(elems[:len(elems)-1], "/")
	name := elems[len(elems)-1]
	f.Import(pkgpath)
	f.Code(func(e *handwriting.Emitter) error {
		return AsFakeStruct(f, e.Prog.Package(pkgpath), name, e.Output, exportedOnly)
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
func AsFakeStruct(f *handwriting.File, info *loader.PackageInfo, name string, o *indent.Output, exportedOnly bool) error {
	target := info.Pkg.Scope().Lookup(name)
	if target == nil {
		return errors.Errorf("%q is not found from package %q", name, info.Pkg.Path())
	}

	named, _ := target.Type().(*types.Named)
	if named == nil {
		return errors.Errorf("%q is not interface", name)
	}
	iface, _ := named.Underlying().(*types.Interface)
	if iface == nil {
		return errors.Errorf("%q is not interface", name)
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
